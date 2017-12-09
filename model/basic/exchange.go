package basic

// exchange

import (
	"encoding/base64"

	"github.com/weihualiu/ewp/router"
	log "github.com/Sirupsen/logrus"
	"github.com/weihualiu/ewp/m"
	"github.com/weihualiu/ewp/utils"
	c "github.com/weihualiu/ewp/model/constant"
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
	buf := make([]byte, 1024)
	_, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		log.Errorf("/user/hello. base64 decode failed!")
		return nil, utils.NewErrorByte([]byte{c.HANDSHAKE_FAILURE})
	}
	
	sessionid, ok := req.GetSessionId()
	if !ok {
		return nil, utils.NewErrorByte([]byte{c.ACCESS_DENIED})
	}
	log.Println("session id:", sessionid)
	
	
	response := m.ResponseNew()
	
	return response, nil
}