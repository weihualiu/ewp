package session

// clientinfo

import (
	"net/http"
	"strings"
)

type ClientInfo struct {
	Id string
	Ad string
	Platform string
	Target string
	Version string
	OSVersion string
}

func ClientInfoNew(header http.Header) *ClientInfo {
	this := new(ClientInfo)
	infostr := header.Get("clientinfo")
	infoarr := strings.Split(infostr, "-")
	if len(infoarr) == 6 {
		this.Id = infoarr[0]
		this.Ad = infoarr[1]
		this.Platform = infoarr[2]
		this.Target = infoarr[3]
		this.Version = infoarr[4]
		this.OSVersion = infoarr[5] // --UA
	}else if len(infoarr) == 5 {
		this.Id = infoarr[0]
		this.Platform = infoarr[1]
		this.Target = infoarr[2]
		this.Version = infoarr[3]
		this.OSVersion = infoarr[4] // --UA
	}else if len(infoarr) == 4 {
		this.Id = infoarr[0]
		this.Platform = infoarr[1]
		this.Target = infoarr[2]
		this.Version = infoarr[3]
	}
	
	return this
}
