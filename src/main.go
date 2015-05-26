package main

import (
	"flag"
	"github.com/qiniu/log"
	"github.com/slene/iploc"
	"path/filepath"
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

	//create task runner
	qlog.GlbTaskRunner = &qlog.QTaskRunner{}
	qlog.GlbTaskRunner.Init()
	go qlog.GlbTaskRunner.Scheduler()
	server := qlog.QLogServer{
		qlog.GlbConf,
	}
	//init database
	qlog.InitDB()
	//load task from database
	loadErr := qlog.LoadTaskFromDB()
	if loadErr != nil {
		log.Error("load task from db error", loadErr.Error())
	}
	//load ip loc info
	ipLocFilePath, _ := filepath.Abs("data/iploc.dat")
	iploc.IpLocInit(ipLocFilePath, true)
	//start to listen
	listenErr := server.Listen()
	if listenErr != nil {
		log.Error("start server error,", listenErr.Error())
	}
}
