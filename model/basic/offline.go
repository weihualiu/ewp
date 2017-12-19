package basic

import (
	"github.com/weihualiu/ewp/sessions"
	c "github.com/weihualiu/ewp/model/constant"
	"strings"
	"bytes"
	"github.com/weihualiu/ewp/utils"
	 log "github.com/Sirupsen/logrus"
)

// offline resources

type Resource struct {

}

func resourceMsg(data []byte, cipher *sessions.CipherState, clientinfo *sessions.ClientInfo) []byte {
	message := messageNew(data)
	hash, _ := message.decode(c.RESOURCE_HASH)
	hashOpt, _ := message.decode(c.RESOURCE_HASH_OPT)
	hashH5, _ := message.decode(c.RESOURCE_H5_HASH)
	if hash == nil && hashOpt == nil && hashH5 == nil {
		return nil
	}

	var ver string
	ver_arr := strings.Split(clientinfo.ResVer, ".")
	if ver_arr[1] == "0" {
		ver = ver_arr[0]
	}else{
		ver = clientinfo.ResVer
	}
	hashRes := bytes.NewBuffer([]byte{})
	if ver == "4" {
		hashRes.Write(resourceUpdate(clientinfo, hash, hashOpt, hashH5))
	}else {

	}
	hashResEncrypt, err := utils.AesEncrypt(hashRes.Bytes(), cipher.ServerKey, cipher.ServerIV)
	if err != nil {
		log.Error("offline hash resource message failed!")
		return nil
	}
	return encodemsg(c.RESOURCE_HASH_RES, hashResEncrypt)
}

func resourceUpdate(client *sessions.ClientInfo, hash, hashOpt, hashH5 []byte) []byte{
	return  nil
}
