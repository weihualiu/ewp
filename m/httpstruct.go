package m

import (
	"net/http"
	"bytes"
	"strings"
	
	log "github.com/Sirupsen/logrus"
)

type Response struct {
	Header map[string]string
	Body *bytes.Buffer
}

func ResponseNew() *Response {
	return &Response{Header:make(map[string]string), Body:bytes.NewBuffer([]byte{})}
}

func (this *Response)Write(p []byte) (n int, err error) {
	return this.Body.Write(p)
}

func (this Response)WriteToResponse(w http.ResponseWriter) {
	for k,v := range this.Header {
		w.Header().Set(k, v)
	}
	log.Println("http response body:", this.Body.Bytes())
	w.Write(this.Body.Bytes())
}

// 封装HTTP请求
type Request struct {
	Req *http.Request
	Body []byte
	BodyParam map[string]string
}

func RequestNew(req *http.Request) *Request {
	return &Request{Req: req}
}

func (this *Request)SetBody(content []byte) {
	this.Body = content
}

func (this *Request)SetBodyParam(content []byte) {
	// string(content)
}

func (this Request)GetParam(key string) string {
	return this.Req.URL.Query().Get(key)
}

func (this Request)GetHeader(key string) string {
	return this.Req.Header.Get(key)
}

func (this Request)GetSessionId() (string, bool) {
	s1 := this.Req.Header.Get("X-Emp-Cookie")
	
	if s1 == "" {
		s2 := this.Req.Header.Get("Cookie")
		if s2 == "" {
			return "", false
		}else{
			return strings.Split(strings.Split(s2,";")[0], "=")[1], true
		}
	}else{
		return strings.Split(s1, "=")[1], true
	}
}
