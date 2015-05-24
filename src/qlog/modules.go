package qlog

import (
	"errors"
	"fmt"
)

//查询日志数据
func QueryLogData(bucket string, dateStr string, params ...string) (records *[]QLogRecord, err error) {
	logStatus, qErr := QueryLogStatus(bucket, dateStr)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query log status error due to, %s", qErr.Error()))
		return
	}

	if logStatus == nil || !logStatus.Done {
		err = errors.New("log data not analysed yet, wait for a few minutes and try again")
		//add task to task queue
		GlbTaskRunner.AddTask(bucket, dateStr)
		return
	} else {
		//query log data

	}

	return
}

//
