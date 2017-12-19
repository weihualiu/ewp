package basic

// handshake
import (
	"encoding/base64"

	log "github.com/Sirupsen/logrus"
	"github.com/weihualiu/ewp/m"
	c "github.com/weihualiu/ewp/model/constant"
	"github.com/weihualiu/ewp/router"
	"github.com/weihualiu/ewp/utils"
	"github.com/weihualiu/ewp/sessions"
	"bytes"
)

func init() {
	r := router.Router{
		CTXType:    router.CTX_BINARY,
		Encrypt:    false,
		CheckValid: false,
		Handler:    handshakeHandler}
	router.Register("/user/handshake", r)
}

func handshakeHandler(data []byte, req *m.Request) (*m.Response, *utils.Error) {
	log.Println("handshake handler")
	buf := make([]byte, len(data))
	_, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		log.Errorf("/user/hello. base64 decode failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	conns := sessions.ConnectionStateInit()
	secopt := sessions.GetSecOptions()

	message := dec_handshake(buf)
	clientHello, decOk := message.decode(c.CLIENT_HELLO)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	client_hello(clientHello, conns, secopt)
	clientInfo := sessions.ClientInfoNew()
	clientInfo.SetAll(req.Req.URL.Query())
	sessions.InsertClientInfo(conns.SessionId, clientInfo)

	clientKey, decOk := message.decode(c.CLIENT_KEY_EXCHANGE)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	client_key_exchange(clientKey, conns, secopt)
	certificate, decOk := message.decode(c.CERTIFICATE)
	if decOk {
		client_certificate(certificate, conns, secopt)
	}
	certificateVerify, decOk := message.decode(c.CERTIFICATE_VERIFY)
	if decOk {
		certificate_verify(certificateVerify, conns, secopt)
	}
	changeCipherSpec, decOk := message.decode(c.CHANGE_CIPHER_SPEC)
	if !decOk {
		log.Error("changeCipherSpec failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	change_cipher_spec(changeCipherSpec, conns, secopt)
	finish, decOk := message.decode(c.FINISHED)
	if !decOk {
		log.Error("finish failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	finished(finish, conns, secopt)
	serverHello := server_hello(conns)
	serverCerticate := server_certificate(conns, secopt)
	serverKeyExchange := server_key_exchange(conns, secopt)
	changeCipherSpec2 := change_cipher_spec2(conns, secopt)
	finished := finished2(conns, secopt)
	cipher := init_cipher_state(conns)
	sessions.InsertConnState(conns.SessionId, conns)
	initContentMsg := init_content(cipher)

	response := m.ResponseNew()
	response.Write(serverHello)
	response.Write(serverCerticate)
	response.Write(serverKeyExchange)
	response.Write(changeCipherSpec2)
	response.Write(finished)
	response.Write(initContentMsg)
	response.Write(resourceMsg(buf, cipher, clientInfo))
	response.SetHeader("Set-Cookie", "_session_id=" + conns.SessionId + "; path=/")
	response.SetHeader("Content-Type", "application/octet-stream")
	return response, nil
}

func client_hello(message []byte, conn *sessions.ConnectionState, secopt *sessions.SecOptions) {
	clientHello := clientHelloNew()
	clientHello.Parse(message)
	version := sessions.SelectVersion(clientHello.Version, secopt.OldestVer)
	sid := sessions.SessionSelectId(clientHello.SessionId)
	cipherSuite := sessions.CipherSuiteSelect(clientHello.CipherS)
	serialNumber := secopt.ServerCert.SerialNumber.Bytes()
	if bytes.Compare(serialNumber, clientHello.SerialNumber) != 0 {
		panic(c.RESET_FULL_HANDSHAKE)
	}
	conn.VerifyD.Update(encodemsg(c.CLIENT_HELLO, message))
	add_seq(int(clientHello.SeqGroup), conn.SessionId)
	conn.Version = version
	conn.Verify = false
	copy(conn.ClientRandom, clientHello.Random)
	conn.SessionId = sid
	conn.CipherS = cipherSuite
}
