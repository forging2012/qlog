package qlog

import (
	"encoding/json"
	"net/http"
	"strings"
)

type RetAPI struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type LogPrepareStatus struct {
	Bucket string `json:"bucket,omitempty"`
	Date   string `json:"date,omitempty"`
	Done   bool   `json:"done,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

func (this *RetAPI) Bytes() []byte {
	data, _ := json.Marshal(this)
	return data
}

//查询日志处理状态
func (this *QLogServer) serveStatusQuery(w http.ResponseWriter, req *http.Request) {
	ret := RetAPI{}
	if req.Method == "GET" {
		ret.Error = "不支持GET方法"
	} else {
		bucket := req.FormValue("bucket")
		date := req.FormValue("date")
		logStatus, err := QueryLogStatus(bucket, date)
		if err != nil {
			ret.Error = "日志状态查询失败!" + err.Error()
		} else {
			if logStatus != nil {
				if logStatus.Done {
					ret.Data = LogPrepareStatus{Done: true, Msg: "日志预处理已完成!"}
				} else {
					ret.Data = LogPrepareStatus{Done: false, Msg: "日志处理中，请稍候..."}
				}
			} else {
				ret.Error = "日志预处理任务不存在，请触发预处理。"
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(ret.Bytes())
}

//预处理日志
func (this *QLogServer) servePrepare(w http.ResponseWriter, req *http.Request) {
	ret := RetAPI{}
	if req.Method == "GET" {
		ret.Error = "不支持GET方法"
	} else {
		bucket := req.FormValue("bucket")
		date := req.FormValue("date")
		logStatus, err := QueryLogStatus(bucket, date)
		if err != nil {
			ret.Error = "日志状态查询失败!" + err.Error()
		} else {
			if logStatus != nil && logStatus.Done {
				ret.Data = LogPrepareStatus{Done: true, Msg: "日志预处理已完成!"}
			} else {
				GlbTaskRunner.AddTask(bucket, date)
				ret.Data = LogPrepareStatus{Done: false, Msg: "日志处理中，请稍候..."}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(ret.Bytes())
}

//预处理中日志查询
func (this *QLogServer) servePrepareQuery(w http.ResponseWriter, req *http.Request) {
	ret := RetAPI{}
	if req.Method == "GET" {
		ret.Error = "不支持GET方法"
	} else {
		logPrepareStatus, err := QueryLogPrepare()
		if err != nil {
			ret.Error = "日志状态查询失败!" + err.Error()
		} else {
			if len(logPrepareStatus) > 0 {
				ret.Data = logPrepareStatus
			} else {
				ret.Data = "日志预处理任务已全部完成!"
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(ret.Bytes())
}

//添加日志分析配置
func (this *QLogServer) serveSettingsAdd(w http.ResponseWriter, req *http.Request) {
	ret := RetAPI{}
	if req.Method == "GET" {
		ret.Error = "不支持GET方法"
	} else {
		bucket := req.FormValue("bucket")
		saveBucket := req.FormValue("saveBucket")
		saveBucketDomain := req.FormValue("saveBucketDomain")
		isSaveBucketPrivateVal := req.FormValue("isSaveBucketPrivate")
		if strings.TrimSpace(bucket) == "" {
			ret.Error = "请填写日志源空间!"
		} else if strings.TrimSpace(saveBucket) == "" {
			ret.Error = "请填写日志存储空间!"
		} else if strings.TrimSpace(saveBucketDomain) == "" {
			ret.Error = "请填写日志存储空间的域名!"
		} else if strings.TrimSpace(isSaveBucketPrivateVal) == "" {
			ret.Error = "请选择日志存储空间的类型!"
		} else {
			isSaveBucketPrivate := false
			if isSaveBucketPrivateVal == "1" {
				isSaveBucketPrivate = true
			}
			err := WriteLogSyncSettings(bucket, saveBucket, saveBucketDomain, isSaveBucketPrivate)
			if err != nil {
				ret.Error = "添加发生错误!" + err.Error()
			} else {
				ret.Data = "添加成功!"
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(ret.Bytes())
}

//删除日志分析配置
func (this *QLogServer) serveSettingsDelete(w http.ResponseWriter, req *http.Request) {
	ret := RetAPI{}
	if req.Method == "GET" {
		ret.Error = "不支持GET方法"
	} else {
		bucket := req.FormValue("bucket")
		err := DeleteLogSyncSettings(bucket)
		if err != nil {
			ret.Error = "删除失败!" + err.Error()
		} else {
			ret.Data = "删除成功!"
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(ret.Bytes())
}
