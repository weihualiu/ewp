package sessions

// clientinfo

import (
	//"net/http"
	"strings"
)

type ClientInfo struct {
	Id string
	Ad string
	Platform string
	Target string
	Version string
	OSVersion string
	DeviceId string
	OtaVer string    // Ex: ND-UMP-3.0.0-080901
}

func ClientInfoNew(str string) *ClientInfo {
	this := new(ClientInfo)
	infoarr := strings.Split(str, "-")
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
