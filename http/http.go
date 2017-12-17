package http

import (
	"net/http"
	g "github.com/weihualiu/ewp/conf"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	_ "github.com/weihualiu/ewp/model"
	"github.com/weihualiu/ewp/sessions"
	"github.com/weihualiu/ewp/router"
	"github.com/weihualiu/ewp/m"
)


func Start() {
	sessions.Init()
	
	addr := g.Config().Http.Addr
	http.HandleFunc("/", ParseRequestHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}


func ParseRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("http request:", r)
	// r.URL.Path
	rt, err := router.Get(r.URL.Path)
	if err != nil {
		w.Write([]byte("not found router path!"))
		return
	}
	
	body, err := ioutil.ReadAll(r.Body)
	
	request := m.RequestNew(r)
	request.Body = body
	
	if rt.CheckValid {
		// 校验会话
		if !validSession(r.Header) {
			w.Write([]byte("session valid!"))
			return
		}
	}
	
	if rt.Encrypt {
		// 执行解密操作
		request.Body = decrypt(body)
	}
	
	defer func() {
		// 处理handler panic的情况
		
	}()
	
	resp_data, e := rt.Handler(request.Body, request)
	if e != nil {
		log.Println(e.Error())
		w.Write(e.Error())
		return
	}
	
	resp_data.WriteToResponse(w)
	
}

func validSession(header http.Header) bool {
	s := header.Get("X-Emp-Cookie")
	if s == "" {
		s = header.Get("Cookie")
	}
	if s == "" {
		return false
	}
	if sessions.GetSession(s) == nil {
		return false
	}
	return true
}

func decrypt(content []byte) []byte {
	return nil
}