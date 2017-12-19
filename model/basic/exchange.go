package basic

// exchange

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"

	log "github.com/Sirupsen/logrus"
	"github.com/weihualiu/ewp/m"
	c "github.com/weihualiu/ewp/model/constant"
	"github.com/weihualiu/ewp/router"
	"github.com/weihualiu/ewp/sessions"
	"github.com/weihualiu/ewp/utils"
)

// 绑定handler到http上
func init() {
	r := router.Router{
		CTXType:    router.CTX_BINARY,
		Encrypt:    false,
		CheckValid: false,
		Handler:    exchangeHandler}
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
	session.Client.SetAll(req.Req.URL.Query())
	conn := session.Conn
	// 更新服务端随机数、预主密钥、主密钥到connectionstate中
	client_key_exchange(clientKey, conn, secopt)
	// 根据配置文件确定是否验证客户端证书
	certificate, decOk := message.decode(c.CERTIFICATE)
	if decOk {
		client_certificate(certificate, conn, secopt)
	}
	// 验证客户端的签名是否正确
	certificateVerify, decOk := message.decode(c.CERTIFICATE_VERIFY)
	if !decOk {
		log.Error("certificateVerify failed!")
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	certificate_verify(certificateVerify, conn, secopt)
	//
	changeCipherSpec, decOk := message.decode(c.CHANGE_CIPHER_SPEC)
	if !decOk {
		log.Error("changeCipherSpec failed!")
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	change_cipher_spec(changeCipherSpec, conn, secopt)

	finish, decOk := message.decode(c.FINISHED)
	if !decOk {
		log.Error("finish failed!")
		return nil, utils.NewErrorByte([]byte{c.UNEXPECTED_MESSAGE})
	}
	finished(finish, conn, secopt)

	serverKeyExchangeBuf := server_key_exchange(conn, secopt)
	serverChangeCipherSpecBuf := change_cipher_spec2(conn, secopt)
	serverFinishedBuf := finished2(conn, secopt)

	cipher := init_cipher_state(conn)
	initContentMsg := init_content(cipher)

	response := m.ResponseNew()
	response.Write(serverKeyExchangeBuf)
	response.Write(serverChangeCipherSpecBuf)
	response.Write(serverFinishedBuf)
	response.Write(initContentMsg)
	response.Write(resourceMsg(buf,cipher, session.Client))
	response.SetHeader("Content-Type", "application/octet-stream")
	return response, nil
}

func change_cipher_spec(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	message := encodemsg(c.CHANGE_CIPHER_SPEC, data)
	conn.VerifyD.Update(message)
	conn.ClientCertificate = data
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

func init_content(cipher *sessions.CipherState) []byte {
	data := bytes.NewBuffer([]byte{})
	data.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>
<content>
    <head>
        <script type="text/x-lua" src="RYTL.lua"></script>
        <script type="text/x-lua">
          <![CDATA[
             function alert_callback()
             	this:setPhysicalkeyListeners();
             end;
             local err_msg = "first";
             window:alert(err_msg, "确定", alert_callback);
            ]]>
        </script>
    </head></content>`))

	encrypt, err := utils.AesEncrypt(data.Bytes(), cipher.ServerKey, cipher.ServerIV)
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
	cipher.ClientIV = keyBlock[keylen : keylen+ivsize]
	cipher.ClientMac = keyBlock[keylen+ivsize : sumlen]
	cipher.ServerKey = keyBlock[sumlen : sumlen+keylen]
	cipher.ServerIV = keyBlock[sumlen+keylen : wantedlen-hashsize]
	cipher.ServerMac = keyBlock[wantedlen-hashsize : wantedlen]

	return cipher
}
