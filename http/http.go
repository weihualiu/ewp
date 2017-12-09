package http

import (
	"net/http"
	g "github.com/weihualiu/ewp/conf"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	_ "github.com/weihualiu/ewp/model"
	_ "github.com/weihualiu/ewp/session"
	"github.com/weihualiu/ewp/router"
	"github.com/weihualiu/ewp/m"
)


func Start() {
	addr := g.Config().Http.Addr
	http.HandleFunc("/", ParseRequestHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}


func ParseRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("http request:", r)
	// r.URL.Path
	rt, err := router.Get(r.URL.Path)
	if err != nil {
		return
	}
	
	body, err := ioutil.ReadAll(r.Body)

	if rt.CheckValid {
		// 校验会话
	}
	
	if rt.Encrypt {
		// 执行解密操作
	}
	
	resp_data, e := rt.Handler(body, m.RequestNew(r))
	if e != nil {
		log.Println(e.Error())
		w.Write(e.Error())
		return
	}
	
	resp_data.WriteToResponse(w)
	
}