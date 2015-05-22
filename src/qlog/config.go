package qlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type QLogConfig struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	ListenHost string `json:"listen_host"`
	ListenPort int    `json:"listen_port"`
	DBServer   string `json:"db_server"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPass     string `json:"db_pass"`
}

func (this *QLogConfig) LoadFromFile(confFile string) error {
	cfgFp, openErr := os.Open(confFile)
	if openErr != nil {
		return errors.New(fmt.Sprintf("open config file error, %s", openErr.Error()))
	}
	defer cfgFp.Close()

	jDecoder := json.NewDecoder(cfgFp)
	decErr := jDecoder.Decode(this)
	if decErr != nil {
		return errors.New(fmt.Sprintf("parse config file error, %s", decErr.Error()))
	}
	return nil
}
