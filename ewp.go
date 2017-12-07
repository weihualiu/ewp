package main

import (
	"flag"
	g "github.com/weihualiu/ewp/conf"
	"github.com/weihualiu/ewp/session"
	"github.com/weihualiu/ewp/http"
	_ "net/http/pprof"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()
	
	g.ParseConfig(*cfg)
	if g.Config().Debug {
		g.InitLog("debug")
	} else {
		g.InitLog("info")
	}
	
	session.Init()
	
	//utils.TimerStartService()
	
	go http.Start()
	
	select {}
}