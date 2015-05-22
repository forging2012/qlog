package main

import (
	"flag"
	"log"
	"qlog"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "config file to start app")
	flag.Parse()

	if configFile == "" {
		log.Println("config file not found")
		return
	}

	cnf := &qlog.QLogConfig{}
	cnfErr := cnf.LoadFromFile(configFile)
	if cnfErr != nil {
		log.Println("load config file error", cnfErr.Error())
		return
	}

	server := qlog.QLogServer{
		cnf,
	}
	listenErr := server.Listen()
	if listenErr != nil {
		log.Println("start server error,", listenErr.Error())
	}
}
