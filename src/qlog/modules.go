package qlog

import (
	"errors"
	"fmt"
)

//查询日志数据
func PrepareLogData(bucket string, date string) (msg string, err error) {
	logStatus, qErr := QueryLogStatus(bucket, date)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("日志状态查询失败, %s", qErr.Error()))
		return
	}

	if logStatus == nil || !logStatus.Done {
		msg = "日志处理中, 请稍候..."
		//add task to task queue
		GlbTaskRunner.AddTask(bucket, date)
		return
	} else {
		//
		msg = "日志预处理完成!"
	}

	return
}

//
