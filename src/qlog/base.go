package qlog

import (
	"strconv"
	"strings"
	"time"
)

type QLog struct {
	ReqIp string
	//@deprecated
	BrowseId string
	//@deprecated
	BrowseName string
	ReqTime    time.Time
	ReqMethod  string
	ReqPath    string
	ReqProto   string

	StatusCode int
	//include header and body
	TotalBytes int
	Referer    string
	UserAgent  string
	Host       string
	Version    string
}

//V1版本日志行解析
func (this *QLog) Parse(logLine string) (ok bool) {
	cnt := len(logLine)
	//delimiter is [ ] "
	logItems := make([]string, 0)
	delimHit := false
	delimIndex := 0
	for i := 0; i < cnt; i++ {
		ch := logLine[i]
		if ch == '[' {
			delimHit = true
			delimIndex = i
		} else if ch == ']' {
			delimHit = false
		} else if ch == '"' {
			delimHit = !delimHit
		} else if ch == ' ' {
			if delimHit {
				continue
			} else {
				//append
				item := logLine[delimIndex:i]
				logItems = append(logItems, item)
				delimIndex = i + 1
			}
		}
	}
	logItems = append(logItems, logLine[delimIndex:])

	itemCnt := len(logItems)
	if logItems[itemCnt-1] != "V1" || itemCnt != 11 {
		ok = false
	}

	//parse
	this.ReqIp = logItems[0]
	this.BrowseId = logItems[1]
	this.BrowseName = logItems[2]
	dateTimeStr := Trim(logItems[3], "[", "]")
	this.ReqTime, _ = ParseDateTime(dateTimeStr)
	reqStr := Trim(logItems[4], "\"", "\"")
	reqParts := strings.Split(reqStr, " ")
	this.ReqMethod = reqParts[0]
	this.ReqPath = reqParts[1]
	this.ReqProto = reqParts[2]
	this.StatusCode, _ = strconv.Atoi(logItems[5])
	this.TotalBytes, _ = strconv.Atoi(logItems[6])
	this.Referer = Trim(logItems[7], "\"", "\"")
	this.UserAgent = Trim(logItems[8], "\"", "\"")
	this.Host = Trim(logItems[9], "\"", "\"")
	this.Version = logItems[10]
	ok = true
	return
}
