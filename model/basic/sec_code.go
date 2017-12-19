package basic

import (
	"bytes"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"

	//log "github.com/Sirupsen/logrus"
	"github.com/weihualiu/ewp/sessions"
	"github.com/weihualiu/ewp/utils"
	mstring "github.com/weihualiu/toolkit/string"
	g "github.com/weihualiu/ewp/conf"
	plug "github.com/weihualiu/ewp/model/plugins"
	c "github.com/weihualiu/ewp/model/constant"
)

type messages struct {
	msgblock map[byte][]byte
}

func messageNew(data []byte) *messages {
	m := new(messages)
	m.msgblock = make(map[byte][]byte, 100)
	m.parse(data)
	return m
}

// 信道握手报文组包
func encodemsg(typeval byte, data []byte) []byte {
	m := bytes.NewBuffer([]byte{typeval})
	m.Write(mstring.IntToBytes(len(data)))
	m.Write(data)
	return m.Bytes()
}

// 信道握手报文解包
func (this messages) decode(typeval byte) ([]byte, bool) {
	flag := false
	var content []byte
	for k, v := range this.msgblock {
		if k == typeval {
			content = v
			flag = true
			break
		}
	}
	return content, flag

}

func (this *messages) parse(data []byte) {
	end := 0
	for end < len(data) {
		if len(data) < end+5 {
			break
		}
		tval := byte(data[end])
		lenBuff := int(mstring.BytesToUInt32(data[end+1 : end+5]))
		if end+5+lenBuff > len(data) {
			break
		}
		this.msgblock[tval] = data[end+5 : end+5+lenBuff]
		end += (lenBuff + 5)
	}

}

type clientHello struct {
	Version      sessions.SecurityVersion
	Random       []byte
	SessionId    string
	SeqGroup     uint8
	CipherS      []sessions.CipherSuite
	SerialNumber []byte
}

func clientHelloNew() *clientHello {
	return new(clientHello)
}

func (this *clientHello) Parse(msg []byte) {
	this.Version = sessions.SecurityVersion{uint8(msg[0]), uint8(msg[1])}
	this.Random = msg[2:34]
	if this.Version.ToInt() >= 104 {
		this.SeqGroup = uint8(msg[34])
		cs_len := mstring.BytesToUInt16(msg[35:37])
		// many ciphersuite
		this.CipherS = sessions.ParseCipherSuites(msg[37 : int(cs_len)+37])
		//this.CipherS = msg[37: int(cs_len) + 37]
		sn_len := int(msg[int(cs_len)+37])
		this.SerialNumber = msg[int(cs_len)+38 : int(cs_len)+38+sn_len]
	} else {
		sid_len := int(msg[34])
		this.SessionId = mstring.BytesToString(msg[35 : sid_len+35])
		cs_len := mstring.BytesToUInt16(msg[sid_len+35 : sid_len+37])
		this.CipherS = sessions.ParseCipherSuites(msg[sid_len+37 : int(cs_len)+sid_len+37])
		sn_len := int(msg[int(cs_len)+sid_len+37])
		this.SerialNumber = msg[int(cs_len)+sid_len+38 : int(cs_len)+sid_len+38+sn_len]
	}

}

func dec_cke(data []byte) (premaster, serverRandom, extension []byte, reflag bool) {
	if len(data) < 81 {
		return nil, nil, nil, false
	}
	premaster = data[0:48]
	serverRandom = data[48:80]
	elen := int(data[80])
	extension = data[81 : elen+81]
	reflag = true
	return
}

func enc_pms(version sessions.SecurityVersion, rand []byte) []byte {
	b := bytes.NewBuffer([]byte{})
	b.Write([]byte{version.Major})
	b.Write([]byte{version.Minor})
	b.Write(rand)

	return b.Bytes()
}

func enc_ske(random, pms, mac, masterSecret []byte) []byte {
	b := bytes.NewBuffer(random)
	b.Write(pms)
	b.Write(mac)

	key := masterSecret[0:32]
	iv := masterSecret[32:48]

	encrypt, err := utils.AesEncrypt(b.Bytes(), key, iv)
	if err != nil {
		return nil
	}
	return encrypt
}

func dec_handshake(data []byte) *messages {
	m := messageNew(data)
	// 对含有数据结构进行判断
	// verify_none  CLIENT_HELLO  CLIENT_KEY_EXCHANGE CHANGE_CIPHER_SPEC FINISHED
	// verify_peer  CLIENT_HELLO CLIENT_KEY_EXCHANGE CERTIFICATE CERTIFICATE_VERIFY CHANGE_CIPHER_SPEC FINISHED
	// 对信道版本进行判断

	return m
}

func add_seq(seqgroup int, sessionid string) {

}

func client_key_exchange(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {

	master_secret1 := func(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) ([]byte, []byte, []byte, error) {
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

	var b1, b2, b3 []byte
	if g.Config().RemoteExchange {
		b1, b2, b3, _ = plug.ExchangeCallback(conn.SessionId, conn.CipherS.Bytes(), data, conn.ClientRandom)
	} else {
		b1, b2, b3, _ = master_secret1(data, conn, secopt)
	}

	message := encodemsg(c.CLIENT_KEY_EXCHANGE, data)
	conn.VerifyD.Update(message)

	conn.ClientPMS = b1
	conn.MasterSecret1 = b2
	if conn.ServerRandom == nil {
		conn.ServerRandom = b3
	}

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

func finished(data []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	// 验证是否和客户端上送的计算后的PRF值一致
	finished_verify()
	message := encodemsg(c.FINISHED, data)
	conn.VerifyD.Update(message)
	conn.ClientCertificate = data
}

func server_hello(connstate *sessions.ConnectionState) []byte {
	sid_len := len(connstate.SessionId)
	data := make([]byte, 1+1+32+1+sid_len+2)
	data[0] = connstate.Version.Major
	data[1] = connstate.Version.Minor
	copy(data[2:34], utils.Random())
	data[34] = uint8(sid_len)
	copy(data[35:sid_len+35], []byte(connstate.SessionId))
	copy(data[sid_len+35:], connstate.CipherS.Bytes())

	data1 := encodemsg(c.SERVER_HELLO, data)
	connstate.VerifyD.Update(data1)

	return data1
}

// 获取证书ASN1 CERT内容
func server_certificate(connstate *sessions.ConnectionState, sec *sessions.SecOptions) []byte {
	data := encodemsg(c.CERTIFICATE, sec.ServerCert.Raw)
	connstate.VerifyD.Update(data)

	return data
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