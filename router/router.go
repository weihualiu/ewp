package router

import (
	"errors"
	"github.com/weihualiu/ewp/utils"
	"github.com/weihualiu/ewp/m"
)

var (
	CTX_BINARY = byte(0x01)
	CTX_JSON = byte(0x02)
	CTX_XML = byte(0x03)
	CTX_URL = byte(0x04)
)

type Router struct {
	CTXType byte  // 内容类型 二进制 JSON XML URL串
	Encrypt bool
	CheckValid bool  //session valid
	Handler func([]byte, *m.Request)(*m.Response, *utils.Error)
	
}

var routeres map[string]Router

func newRouters() {
	if routeres == nil {
		routeres = make(map[string]Router)
	}
}

// 注册路径扩展版
func RegisterExt(reqpath string, encrypt bool, check bool, handler func([]byte, *m.Request)(*m.Response, *utils.Error)) error {
	newRouters()
	router := Router{
		CTXType : CTX_URL,
		Encrypt : encrypt,
		CheckValid : check,
		Handler : handler}
	routeres[reqpath] = router
	return nil
}

// 注册路径
func Register(reqpath string, router Router) error {
	newRouters()
	routeres[reqpath] = router
	return nil
}

// 根据路径获取Router
func Get(reqpath string) (r Router, err error) {
	r, ok := routeres[reqpath]
	if !ok {
		return r, errors.New("not found router!")
	}
	return r, nil
}

