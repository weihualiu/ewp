package m

import (
	"net/http"
)

type Response struct {
	Header map[string]string
	Body []byte
}

func (this Response)WriteToResponse(w http.ResponseWriter) {
	for k,v := rang this.Header {
		w.Header().Set(k, this.Header[k])
	}
	w.Write(this.Body)
}