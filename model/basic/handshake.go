package basic

// handshake
import (
	"encoding/base64"

	"github.com/weihualiu/ewp/utils"
	"github.com/weihualiu/ewp/m"
	c "github.com/weihualiu/ewp/model/constant"
	"github.com/weihualiu/ewp/router"
	log "github.com/Sirupsen/logrus"
)

func init() {
	r := router.Router{
		CTXType : router.CTX_BINARY,
		Encrypt : false,
		CheckValid : false,
		Handler : handshakeHandler}
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
	
	response := m.ResponseNew()
	
	return response, nil
}