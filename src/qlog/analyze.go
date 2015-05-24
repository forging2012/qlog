package qlog

import (
	"bufio"
	"compress/gzip"
	"errors"
	"github.com/qiniu/log"
	"os"
	"strings"
)

func ParseLogContent(bucket string, date string, paths []string) (err error) {
	for _, path := range paths {
		lfp, openErr := os.Open(path)
		if openErr != nil {
			err = errors.New("open log file failed due to, " + openErr.Error())
			return
		}
		gzReader, gErr := gzip.NewReader(lfp)
		if gErr != nil {
			err = errors.New("open gz file failed due to, " + gErr.Error())
			return
		}
		bReader := bufio.NewScanner(gzReader)
		bReader.Split(bufio.ScanLines)
		for bReader.Scan() {
			line := strings.TrimSpace(bReader.Text())
			pLog, pErr := ParseLogLine(line)
			if pErr != nil {
				log.Warn("invalid line `", line, "' in file", path)
			} else {
				id := sha1Hash(bucket + ":" + date + ":" + line)
				err := WriteQLogRecord(id, bucket, date, pLog.ReqIp, pLog.ReqTime, pLog.ReqMethod, pLog.ReqProto, pLog.StatusCode,
					pLog.TotalBytes, pLog.Referer, pLog.UserAgent, pLog.Host, pLog.Version)
				if err != nil {
					log.Warn("write log record to db failed due to," + err.Error())
				}
			}
		}
		lfp.Close()
	}
	return
}
