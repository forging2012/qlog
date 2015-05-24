package main

import (
	"flag"
	"github.com/qiniu/log"
	"qlog"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "config file to start app")
	flag.Parse()

	if configFile == "" {
		log.Error("config file not found")
		return
	}

	qlog.GlbConf = &qlog.QLogConfig{}
	cnfErr := qlog.GlbConf.LoadFromFile(configFile)
	if cnfErr != nil {
		log.Error("load config file error", cnfErr.Error())
		return
	}

	//
	qlog.GlbTaskRunner = &qlog.QTaskRunner{}
	qlog.GlbTaskRunner.Init()
	go qlog.GlbTaskRunner.Scheduler()
	server := qlog.QLogServer{
		qlog.GlbConf,
	}
	qlog.InitDB()
	listenErr := server.Listen()
	if listenErr != nil {
		log.Error("start server error,", listenErr.Error())
	}
}
