package sessions

// clientinfo

import (
	"strings"
	"net/url"
)

type ClientInfo struct {
	Id        string
	Ad        string
	Platform  string
	Target    string
	Version   string
	OSVersion string
	DeviceId  string
	OtaVer    string // Ex: ND-UMP-3.0.0-080901
	Resolution string
	ResVer string
	AppName string
}

func ClientInfoNew() *ClientInfo{
	return new(ClientInfo)
}
func (this *ClientInfo)SetAll(values url.Values) {
	this.SetClient(values.Get("clientinfo"))
	this.SetDeviceId(values.Get("diviceId"))
	this.SetOtaVer(values.Get("ota_version"))
	this.SetResolution(values.Get("resolution"))
	if this.Platform == "" {
		this.SetPlatform(values.Get("platform"))
	}
	this.SetResVer(values.Get("res_ver"))
	this.SetAppName(values.Get("app"))
}

func (this *ClientInfo)SetAppName(value string) {
	this.AppName = value
}

func (this *ClientInfo)SetResVer(value string) {
	this.ResVer = value
}

func (this *ClientInfo)SetPlatform(value string) {
	this.Platform = value
}

func (this *ClientInfo)SetResolution(value string) {
	this.Resolution = value
}

func (this *ClientInfo)SetOtaVer(value string) {
	this.OtaVer = value
}

func (this *ClientInfo)SetDeviceId(value string) {
	this.DeviceId = value
}

func (this *ClientInfo)SetClient(value string){
	infoarr := strings.Split(value, "-")
	if len(infoarr) == 6 {
		this.Id = infoarr[0]
		this.Ad = infoarr[1]
		this.Platform = infoarr[2]
		this.Target = infoarr[3]
		this.Version = infoarr[4]
		this.OSVersion = infoarr[5] // --UA
	} else if len(infoarr) == 5 {
		this.Id = infoarr[0]
		this.Platform = infoarr[1]
		this.Target = infoarr[2]
		this.Version = infoarr[3]
		this.OSVersion = infoarr[4] // --UA
	} else if len(infoarr) == 4 {
		this.Id = infoarr[0]
		this.Platform = infoarr[1]
		this.Target = infoarr[2]
		this.Version = infoarr[3]
	}
}
