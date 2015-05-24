package qlog

import (
	"errors"
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
func ParseLogLine(logLine string) (log *QLog, err error) {
	logLine = strings.TrimSpace(logLine)
	cnt := len(logLine)
	//delimiter is [ ] "
	delimStatus := map[uint8]bool{
		'[': false,
		'"': false,
	}
	logItems := make([]string, 0)
	delimStartIndex := 0
	for i := 0; i < cnt; i++ {
		ch := logLine[i]
		switch ch {
		case '[':
			_, exists := existOtherHit(delimStatus, ch)
			if !exists {
				delimStartIndex = i
				delimStatus[ch] = true
			}
		case ']':
			och, exists := existOtherHit(delimStatus, ch)
			if exists && och == '[' {
				delimStatus['['] = false
			}
		case '"':
			_, exists := existOtherHit(delimStatus, ch)
			if !exists {
				if delimStatus[ch] {
					delimStatus[ch] = false
				} else {
					delimStartIndex = i
					delimStatus[ch] = true
				}
			}
		case ' ':
			_, exists := existOtherHit(delimStatus, ch)
			if !exists {
				item := logLine[delimStartIndex:i]
				logItems = append(logItems, item)
				delimStartIndex = i + 1
			}
		default:
			continue
		}
	}
	logItems = append(logItems, logLine[delimStartIndex:])
	itemCnt := len(logItems)
	if logItems[itemCnt-1] != "V1" || itemCnt != 11 {
		err = errors.New("invalid log line")
		return
	}

	log = &QLog{}
	//parse
	log.ReqIp = logItems[0]
	log.BrowseId = logItems[1]
	log.BrowseName = logItems[2]
	dateTimeStr := Trim(logItems[3], "[", "]")
	log.ReqTime, _ = ParseDateTime(dateTimeStr)
	reqStr := Trim(logItems[4], "\"", "\"")
	reqParts := strings.Split(reqStr, " ")
	log.ReqMethod = reqParts[0]
	log.ReqPath = reqParts[1]
	log.ReqProto = reqParts[2]
	log.StatusCode, _ = strconv.Atoi(logItems[5])
	log.TotalBytes, _ = strconv.Atoi(logItems[6])
	log.Referer = Trim(logItems[7], "\"", "\"")
	log.UserAgent = Trim(logItems[8], "\"", "\"")
	log.Host = Trim(logItems[9], "\"", "\"")
	log.Version = logItems[10]
	return
}

func existOtherHit(status map[uint8]bool, ch uint8) (och uint8, exists bool) {
	for c, e := range status {
		if c != ch && e {
			och = c
			exists = e
		}
	}
	return
}
