package basic

// /usr/hello

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	
	"github.com/weihualiu/ewp/router"
	"github.com/weihualiu/ewp/utils"
	mstring "github.com/weihualiu/toolkit/string"
	c "github.com/weihualiu/ewp/model/constant"
	"github.com/weihualiu/ewp/m"
	"github.com/weihualiu/ewp/sessions"
)

// 绑定handler到http上
func init() {
	r := router.Router{
		CTXType : router.CTX_BINARY,
		Encrypt : false,
		CheckValid : false,
		Handler : helloHandler}
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
	sessions.InsertClientInfo(connstate.SessionId, sessions.ClientInfoNew(req.GetParam("clientinfo")))
	// 设置Cookie值
	response.Header["Set-Cookie"] = "_session_id=" + connstate.SessionId + "; path=/"
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
	if int(lenBuff) + 1 + 4 > len(data) {
		return false, nil
	}
	
	return true,data[5:lenBuff+5]
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

func server_hello(connstate *sessions.ConnectionState) []byte {
	sid_len := len(connstate.SessionId)
	// <<?BYTE(Major), ?BYTE(Minor),Random:32/binary,?BYTE(SID_length),Session_ID/binary,Cipher_suite/binary>>
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


