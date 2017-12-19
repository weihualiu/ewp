package basic

// /usr/hello

import (
	log "github.com/Sirupsen/logrus"
	"net/http"

	"github.com/weihualiu/ewp/m"
	c "github.com/weihualiu/ewp/model/constant"
	"github.com/weihualiu/ewp/router"
	"github.com/weihualiu/ewp/sessions"
	"github.com/weihualiu/ewp/utils"
	mstring "github.com/weihualiu/toolkit/string"
)

// 绑定handler到http上
func init() {
	r := router.Router{
		CTXType:    router.CTX_BINARY,
		Encrypt:    false,
		CheckValid: false,
		Handler:    helloHandler}
	router.Register("/user/hello", r)
}

func helloHandler(data []byte, req *m.Request) (*m.Response, *utils.Error) {
	buf, err := utils.Base64Decode(data)
	if err != nil {
		log.Errorf("/user/hello. base64 decode failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	// type(1) + len(4) + data(N)
	buf2, decOk := messageNew(buf).decode(c.CLIENT_HELLO)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	log.Println("buf2:", buf2)

	sec := sessions.GetSecOptions()
	connstate := sessions.ConnectionStateInit()

	client_hello_1(buf2, connstate, sec)
	log.Println("ConnectionState1")
	hello := server_hello(connstate)
	cert := server_certificate(connstate, sec)
	certreq := certificate_request(connstate, sec)
	sessions.InsertConnState(connstate.SessionId, connstate)

	response := m.ResponseNew()
	// HTTP HEADER增加第一次请求标识
	response.Header["is_first"] = req.GetParam("is_first")
	// 保存客户端信息到会话中
	clientInfo := sessions.ClientInfoNew()
	clientInfo.SetClient(req.GetParam("clientinfo"))
	sessions.InsertClientInfo(connstate.SessionId,clientInfo)
	// 设置Cookie值
	response.SetHeader("Set-Cookie", "_session_id=" + connstate.SessionId + "; path=/")
	response.SetHeader("X-Emp-CipherExpiry", "300000")
	// 生成失效时间
	// 组装最终报文并返回
	response.Write(hello)
	response.Write(cert)
	response.Write(certreq)

	return response, nil
}

// 获取message data
func decHello(data []byte) (bool, []byte) {
	if data[0] != c.CLIENT_HELLO {
		return false, nil
	}
	lenBuff := mstring.BytesToUInt32(data[1:5])
	if int(lenBuff)+1+4 > len(data) {
		return false, nil
	}

	return true, data[5 : lenBuff+5]
}

func client_hello_1(msgdata []byte, connstate *sessions.ConnectionState, sec *sessions.SecOptions) {
	chello := clientHelloNew()
	chello.Parse(msgdata)
	version := sessions.SelectVersion(chello.Version, sec.OldestVer)
	sid := sessions.NewSession()
	// 判断客户端上送的加密套件是否在服务端支持的范围内
	// 同时选取一个加密套件
	cipherSuite := sessions.CipherSuiteSelect(chello.CipherS)
	//
	connstate.Version = version
	connstate.Verify = false
	connstate.ClientRandom = chello.Random
	connstate.SessionId = sid
	connstate.CipherS = cipherSuite
	connstate.VerifyD.Update(encodemsg(c.CLIENT_HELLO, msgdata))

}

func certificate_request(connstate *sessions.ConnectionState, sec *sessions.SecOptions) []byte {
	data := encodemsg(c.CERTIFICATE_REQUEST, sec.ServerCert.Raw)
	connstate.VerifyD.Update(data)

	return data
}

// 从HTTP HEADER获取Session
func getSIDFroRequestHeader(header http.Header) []byte {
	return nil
}

func select_cipher_suite() {

}
