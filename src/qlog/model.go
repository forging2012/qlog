package qlog

import (
	"time"
)

type QLogRecord struct {
	Id         int
	ReqId      string
	ReqTime    time.Time
	ReqMethod  string
	ReqPath    string
	ReqProto   string
	StatusCode int
	TotalBytes int
	Referer    string
	UserAgent  string
	Host       string
	Version    string
}

type LogSyncStatus struct {
	Bucket     string
	Date       string
	StatusCode int
	Status     string
}

//日志同步的配置
type LogSyncSettings struct {
	Bucket              string
	SaveBucket          string
	SaveDomain          string
	IsSaveBucketPrivate bool
}
