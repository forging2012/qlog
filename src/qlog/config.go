package qlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"os"
)

var GlbConf *QLogConfig

type QLogConfig struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	ListenHost string `json:"listen_host"`
	ListenPort int    `json:"listen_port"`
	DBServer   string `json:"db_server"`
	DBPort     int    `json:"db_port"`
	DBName     string `json:"db_name"`
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

func (this *QLogConfig) SQLDataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		this.DBUser, this.DBPass, this.DBServer, this.DBPort, this.DBName)
}

func (this *QLogConfig) Mac() *digest.Mac {
	return &digest.Mac{this.AccessKey, []byte(this.SecretKey)}
}
