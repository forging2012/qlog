package qlog

import (
	"time"
)

//这里定义和数据库表结构相对应的结构体

//日志记录
type QLogRecord struct {
	Id         string
	Bucket     string
	Date       string
	ReqIp      string
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

//日志分析状态
type QLogSyncStatus struct {
	Id     string
	Bucket string
	Date   string
	Done   bool
}

//日志同步的配置
type QLogSyncSettings struct {
	Bucket              string
	SaveBucket          string
	SaveBucketDomain    string
	IsSaveBucketPrivate bool
}
