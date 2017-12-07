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

func helloHandler(data []byte, req *http.Request) ([]byte, *utils.Error) {
	log.Println("hello handler")
	buf := make([]byte, 1024)

	_, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		log.Errorf("/user/hello. base64 decode failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	//buf1 := bytes.TrimRight(buf, string(0x00))
	log.Println("/user/hello body:", buf)
	//1 0 0 0 42 1 4 0 21 0 20 56 179 129 22 229 121 72 221 171 243 194 86 37 186 136 29 187 80 106 57 205 156 49 255 224 175 68 18 10 0 4 0 7 0 6
	// type(1) + len(4) + data(N)
	decOk, buf2 := decHello(buf)
	if !decOk {
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	log.Println("buf2:", buf2)
	
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

