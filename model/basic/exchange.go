package basic

// exchange

import (
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"crypto/hmac"
	"crypto/sha1"
	"bytes"

	"github.com/weihualiu/ewp/router"
	log "github.com/Sirupsen/logrus"
	"github.com/weihualiu/ewp/m"
	"github.com/weihualiu/ewp/utils"
	c "github.com/weihualiu/ewp/model/constant"
	"github.com/weihualiu/ewp/sessions"
	g "github.com/weihualiu/ewp/conf"
	plug "github.com/weihualiu/ewp/model/plugins"
)

// 绑定handler到http上
func init() {
	r := router.Router{
		CTXType : router.CTX_BINARY,
		Encrypt : false,
		CheckValid : false,
		Handler : exchangeHandler}
	router.Register("/user/exchange", r)
}

func exchangeHandler(data []byte, req *m.Request) (*m.Response, *utils.Error) {
	log.Println("exchange handler")
	
	buf, err := utils.Base64Decode(data)
	if err != nil {
		log.Errorf("/user/hello. base64 decode failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	sessionid, ok := req.GetSessionId()
	if !ok {
		return nil, utils.NewErrorByte([]byte{c.ACCESS_DENIED})
	}
	log.Println("session id:", sessionid)
	
	message := messageNew(buf)
	clientKey, decOk := message.decode(c.CLIENT_KEY_EXCHANGE)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	
	secopt := sessions.GetSecOptions()
	session := sessions.GetSession(sessionid)
	session.Client.DeviceId = req.GetParam("diviceId")
	session.Client.OtaVer = req.GetParam("ota_version")
	conn := session.Conn
	// 更新服务端随机数、预主密钥、主密钥到connectionstate中
	client_key_exchange(clientKey, conn, secopt)
	// 根据配置文件确定是否验证客户端证书
	certificate, decOk := message.decode(c.CERTIFICATE)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	client_certificate(certificate, conn, secopt)
	// 验证客户端的签名是否正确
	certificateVerify, decOk := message.decode(c.CERTIFICATE_VERIFY)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	certificate_verify(certificateVerify, conn, secopt)
	// 
	changeCipherSpec, decOk := message.decode(c.CHANGE_CIPHER_SPEC)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	change_cipher_spec(changeCipherSpec, conn, secopt)
	
	finish, decOk := message.decode(c.FINISHED)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	finished(finish, conn, secopt)
	
	//serverHelloBuf := server_hello(conn, secopt)
	//serverCertificateBuf := server_certificate(conn, secopt)
	serverKeyExchangeBuf := server_key_exchange(conn, secopt)
	serverChangeCipherSpecBuf := change_cipher_spec2(conn, secopt)
	serverFinishedBuf := finished2(conn, secopt)
	
	cipher := init_cipher_state(conn)
	initContentMsg := init_content([]byte{}, cipher)
	
	response := m.ResponseNew()
	response.Write(serverKeyExchangeBuf)
	response.Write(serverChangeCipherSpecBuf)
	response.Write(serverFinishedBuf)
	response.Write(initContentMsg)
	
	return response, nil
}

func client_key_exchange(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	
	master_secret1 := func(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) ([]byte,[]byte,[]byte, error) {
		// 私钥解密数据
		rasPrivKey, err := x509.ParsePKCS1PrivateKey(secopt.ServerKey)
		if err != nil {
			return nil, nil, nil, err
		}
		
		plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, rasPrivKey, data)
		preMasterSecret, serverRandom, _, _ := dec_cke(plainText)
		random := bytes.NewBuffer(conn.ClientRandom)
		random.Write(serverRandom)
		masterSecret := sessions.Prf(preMasterSecret, []byte("master secret"), random.Bytes(), 68)
		return preMasterSecret, masterSecret, serverRandom, nil
	}
	
	var b1,b2,b3 []byte
	if g.Config().RemoteExchange {
		b1,b2,b3, _ = plug.ExchangeCallback(conn.SessionId, conn.CipherS.Bytes(), data, conn.ClientRandom)
	}else{
		b1,b2,b3, _ = master_secret1(data, conn, secopt)
	}
	
	message := encodemsg(c.CLIENT_KEY_EXCHANGE, data)
	conn.VerifyD.Update(message)
	
	conn.ClientPMS = b1
	conn.MasterSecret1 = b2
	conn.ServerRandom = b3
}

func client_certificate(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	if conn.Verify {
		// 通过CA验证服务端证书
	}
	message := encodemsg(c.CERTIFICATE, data)
	conn.VerifyD.Update(message)
	conn.ClientCertificate = data
}

func certificate_verify(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	// 签名验证
	signature_verify()
	message := encodemsg(c.CERTIFICATE_VERIFY, data)
	conn.VerifyD.Update(message)
	conn.ClientCertificate = data
}

func change_cipher_spec(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	message := encodemsg(c.CHANGE_CIPHER_SPEC, data)
	conn.VerifyD.Update(message)
	conn.ClientCertificate = data
}

func finished(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	// 验证是否和客户端上送的计算后的PRF值一致
	finished_verify()
	message := encodemsg(c.FINISHED, data)
	conn.VerifyD.Update(message)
	conn.ClientCertificate = data
}

func server_key_exchange(conn *sessions.ConnectionState, secopt *sessions.SecOptions) []byte {
	preMasterSecret := enc_pms(conn.Version, utils.Random2(46))
	random := bytes.NewBuffer(conn.ClientRandom)
	random.Write(conn.ServerRandom)
	masterSecret := sessions.Prf(preMasterSecret, []byte("master secret2"), random.Bytes(), 48)
	nextServerRandom := utils.Random()
	random2 := bytes.NewBuffer(nextServerRandom)
	random2.Write(preMasterSecret)
	hmacSha1 := hmac_sign(random2.Bytes(), conn.MasterSecret1)
	cipher := enc_ske(nextServerRandom, preMasterSecret, hmacSha1, conn.MasterSecret1)
	//SERVER_KEY_EXCHANGE
	message := encodemsg(c.SERVER_KEY_EXCHANGE, cipher)
	conn.VerifyD.Update(message)
	conn.ServerPMS = preMasterSecret
	conn.MasterSecret2 = masterSecret
	
	return message
}

func change_cipher_spec2(conn *sessions.ConnectionState, secopt *sessions.SecOptions) []byte {
	message := encodemsg(c.CHANGE_CIPHER_SPEC, []byte{byte(0x02)})
	conn.VerifyD.Update(message)
	return message
}

func finished2(conn *sessions.ConnectionState, secopt *sessions.SecOptions) []byte {
	//finished_verify_data
	conn.VerifyD.Finish(conn.MasterSecret2)
	message := encodemsg(c.FINISHED, conn.VerifyD.Buffer)
	return message
}


// 签名验证
func signature_verify() error {
	return nil
}

func finished_verify() error {
	
	return nil
}

func hmac_sign(data, masterSecret []byte) []byte {
	mac := hmac.New(sha1.New, masterSecret[48:68])
	mac.Write(data)
	return mac.Sum(nil)
}

func init_content(data []byte, cipher *sessions.CipherState) []byte {
	encrypt, err := utils.AesEncrypt(data, cipher.ServerKey, cipher.ServerIV)
	if err != nil {
		return nil
	}
	
	message := encodemsg(c.INIT_CONTENT, encrypt)
	return message
}

func init_cipher_state(conn *sessions.ConnectionState) *sessions.CipherState {
	hashsize := 20
	keylen := 32
	ivsize := 16
	sumlen := hashsize + keylen + ivsize
	wantedlen := 2 * (hashsize + keylen + ivsize)
	random := bytes.NewBuffer(conn.ClientRandom)
	random.Write(conn.ServerRandom)
	keyBlock := sessions.Prf(conn.MasterSecret2, []byte("key expansion"), random.Bytes(), wantedlen)
	
	cipher := new(sessions.CipherState)
	cipher.Version = conn.Version
	cipher.CipherS = conn.CipherS
	cipher.ClientKey = keyBlock[0:keylen]
	cipher.ClientIV = keyBlock[keylen:keylen+ivsize]
	cipher.ClientMac = keyBlock[keylen+ivsize:sumlen]
	cipher.ServerKey = keyBlock[sumlen:sumlen+keylen]
	cipher.ServerIV = keyBlock[sumlen+keylen:wantedlen-hashsize]
	cipher.ServerMac = keyBlock[wantedlen-hashsize:wantedlen]
	
	return cipher
}