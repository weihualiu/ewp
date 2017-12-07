package basic

// /usr/hello

import (
	"github.com/weihualiu/ewp/router"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"encoding/base64"
	"github.com/weihualiu/ewp/utils"
	mstring "github.com/weihualiu/toolkit/string"
	//"bytes"
	c "github.com/weihualiu/ewp/model/constant"
	"github.com/weihualiu/ewp/m"

	"encoding/pem"
	"crypto/x509"
	"io/ioutil"
	"github.com/weihualiu/ewp/session"
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

func helloHandler(data []byte, req *http.Request) (*m.Response, *utils.Error) {
	buf := make([]byte, 1024)

	_, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		log.Errorf("/user/hello. base64 decode failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	// type(1) + len(4) + data(N)
	decOk, buf2 := decHello(buf)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	log.Println("buf2:", buf2)
	
	sec := session.SecOptionsInit()
	connstate := session.ConnectionStateInit()
	
	client_hello_1(buf2, connstate, sec)
	log.Println("ConnectionState1")
	server_hello(connstate)
	server_certificate(connstate, sec)
	// HTTP HEADER增加第一次请求标识
	// 保存客户端信息到会话中
	// 设置Cookie值
	// 生成失效时间
	// 组装最终报文并返回
	
	return nil, nil
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


func client_hello_1(msgdata []byte, connstate *session.ConnectionState, sec *session.SecOptions) {
	chello := clientHelloNew()
	chello.Parse(msgdata)
	version := select_version(ch.Version, sec.OldestVer)
	sid := sessions.NewSession()
	// 判断客户端上送的加密套件是否在服务端支持的范围内
	// 同时选取一个加密套件
	cipherSuite := nil
	// 
	connstate.Version = version
	connstate.Verify = false
	connstate.ClientRandom = ch.Random
	connstate.SessionId = sid
	connstate.CipherSuite = cipherSuite
	connstate.VerifyD.Update(enc_msg(c.CLIENT_HELLO, msgdata))
	
}

func server_hello(connstate *session.ConnectionState) []byte {
	sid_len := connstate.SessionId
	// <<?BYTE(Major), ?BYTE(Minor),Random:32/binary,?BYTE(SID_length),Session_ID/binary,Cipher_suite/binary>>
	data := make([]byte, 1+1+32+1+sid_len+len(connstate.CipherSuite))
	copy(data, connstate.MajorVer)
	copy(data[1:], connstate.MinorVer)
	copy(data[2:34], random())
	copy(data[34], uint8(sid_len))
	copy(data[35:sid_len+35], connstate.SessionId)
	copy(data[sid_len+35:], connstate.CipherSuite)
	
	connstate.VerifyD.Update(enc_msg(c.SERVER_HELLO, data))
	
	return data
}

func server_certificate(connstate *session.ConnectionState, sec *session.SecOptions) []byte {
	
}

// 选择版本，如果客户端版本低于服务端版本使用客户端版本；否则使用服务端版本
func select_version(clientver []int, oldver []int) {
	return clientver
}

// 从HTTP HEADER获取Session
func getSIDFroRequestHeader(header http.Header) []byte {
	return nil
}


func select_cipher_suite() {
	
}


func random() []byte {
	// 当前时间
	t := mstring.IntToBytes(time.Now().Nanosecond())
	
	c := 28
	b := make([]byte, c)
	rand.Read(b)
	
	r := make([]byte, 32)
	copy(r[0:4], t)
	copy(r[4:], b)
	
	return r
}
