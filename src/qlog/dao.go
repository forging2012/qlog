package qlog

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qiniu/log"
	"time"
)

var (
	glbDB *sql.DB
)

func InitDB() {
	var err error
	glbDB, err = sql.Open("mysql", GlbConf.SQLDataSource())
	if err != nil {
		log.Error("failed to open database due to, %s", err.Error())
		return
	}
	glbDB.SetMaxIdleConns(20)
	glbDB.SetMaxOpenConns(20)
}

//检查日志的状态
func QueryLogStatus(bucket string, date string) (logStatus *QLogSyncStatus, err error) {
	queryStr := "SELECT id,bucket,date,done FROM log_sync_status WHERE bucket=? AND date=?"
	rows, qErr := glbDB.Query(queryStr, bucket, date)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query failed due to, %s", qErr.Error()))
		return
	}

	if rows.Next() {
		logStatus = &QLogSyncStatus{}
		var id string
		var bucket string
		var date string
		var done bool

		rErr := rows.Scan(&id, &bucket, &date, &done)
		if rErr != nil {
			err = errors.New(fmt.Sprintf("read row data failed due to, %s", rErr.Error()))
			return
		}
		logStatus.Id = id
		logStatus.Bucket = bucket
		logStatus.Date = date
		logStatus.Done = done
	}
	return
}

func QueryLogPrepare() (logPrepareStatus []*LogPrepareStatus, err error) {
	queryStr := "SELECT bucket,date,done FROM log_sync_status WHERE done=? ORDER BY bucket,date"
	rows, qErr := glbDB.Query(queryStr, false)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query failed due to, %s", qErr.Error()))
		return
	}

	logPrepareStatus = make([]*LogPrepareStatus, 0)
	for rows.Next() {
		logStatus := &LogPrepareStatus{}
		var bucket string
		var date string
		var done bool

		rErr := rows.Scan(&bucket, &date, &done)
		if rErr != nil {
			err = errors.New(fmt.Sprintf("read row data failed due to, %s", rErr.Error()))
			return
		}

		logStatus.Bucket = bucket
		logStatus.Date = date
		logStatus.Done = done
		logPrepareStatus = append(logPrepareStatus, logStatus)
	}
	return
}

//写入日志的状态
func WriteLogStatus(bucket string, date string, done bool) (err error) {
	stmt, sErr := glbDB.Prepare("INSERT INTO log_sync_status (id,bucket,date,done) " +
		" VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE done=?")
	if sErr != nil {
		err = errors.New(fmt.Sprintf("prepare exec failed due to, %s", sErr.Error()))
		return
	}
	defer stmt.Close()
	id := base64.URLEncoding.EncodeToString([]byte(bucket + ":" + date))
	_, execErr := stmt.Exec(id, bucket, date, done, done)
	if execErr != nil {
		err = errors.New(fmt.Sprintf("failed to insert or update due to, %s", execErr.Error()))
		return
	}
	return
}

//读取同步配置信息
func GetLogSyncSettings(bucket string) (settings *QLogSyncSettings, err error) {
	queryStr := "SELECT bucket,save_bucket,save_bucket_domain,is_save_bucket_private FROM log_sync_settings WHERE bucket=?"
	rows, qErr := glbDB.Query(queryStr, bucket)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query failed due to, %s", qErr.Error()))
		return
	}

	if rows.Next() {
		settings = &QLogSyncSettings{}

		var bucket string
		var saveBucket string
		var saveBucketDomain string
		var isSaveBucketPrivate bool

		rErr := rows.Scan(&bucket, &saveBucket, &saveBucketDomain, &isSaveBucketPrivate)
		if rErr != nil {
			err = errors.New(fmt.Sprintf("read row data failed due to, %s", rErr.Error()))
			return
		}
		settings.Bucket = bucket
		settings.SaveBucket = saveBucket
		settings.SaveBucketDomain = saveBucketDomain
		settings.IsSaveBucketPrivate = isSaveBucketPrivate
	}
	return
}

func GetLogSyncSettingsAll() (settingsAll []*QLogSyncSettings, err error) {
	queryStr := "SELECT bucket,save_bucket,save_bucket_domain,is_save_bucket_private FROM log_sync_settings ORDER BY bucket ASC"
	rows, qErr := glbDB.Query(queryStr)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query failed due to, %s", qErr.Error()))
		return
	}
	settingsAll = make([]*QLogSyncSettings, 0)
	for rows.Next() {
		settings := &QLogSyncSettings{}

		var bucket string
		var saveBucket string
		var saveBucketDomain string
		var isSaveBucketPrivate bool

		rErr := rows.Scan(&bucket, &saveBucket, &saveBucketDomain, &isSaveBucketPrivate)
		if rErr != nil {
			err = errors.New(fmt.Sprintf("read row data failed due to, %s", rErr.Error()))
			return
		}
		settings.Bucket = bucket
		settings.SaveBucket = saveBucket
		settings.SaveBucketDomain = saveBucketDomain
		settings.IsSaveBucketPrivate = isSaveBucketPrivate
		settingsAll = append(settingsAll, settings)
	}
	return
}

//写入同步配置信息
func WriteLogSyncSettings(bucket string, saveBucket string, saveBucketDomain string,
	isSaveBucketPrivate bool) (err error) {
	stmt, sErr := glbDB.Prepare("INSERT INTO log_sync_settings (bucket,save_bucket,save_bucket_domain,is_save_bucket_private) " +
		" VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE save_bucket=?,save_bucket_domain=?,is_save_bucket_private=?")
	if sErr != nil {
		err = errors.New(fmt.Sprintf("prepare exec failed due to, %s", sErr.Error()))
		return
	}
	defer stmt.Close()
	_, execErr := stmt.Exec(bucket, saveBucket, saveBucketDomain, isSaveBucketPrivate, saveBucket, saveBucketDomain, isSaveBucketPrivate)
	if execErr != nil {
		err = errors.New(fmt.Sprintf("failed to insert or update due to, %s", execErr.Error()))
		return
	}
	return
}

//删除同步配置信息
func DeleteLogSyncSettings(bucket string) (err error) {
	stmt, sErr := glbDB.Prepare("DELETE FROM log_sync_settings WHERE bucket=?")
	if sErr != nil {
		err = errors.New(fmt.Sprintf("prepare exec failed due to, %s", sErr.Error()))
		return
	}
	defer stmt.Close()
	_, execErr := stmt.Exec(bucket)
	if execErr != nil {
		err = errors.New(fmt.Sprintf("failed to insert or update due to, %s", execErr.Error()))
		return
	}
	return
}

//获取已配置空间列表
func GetBucketListFromSettings() (buckets []string, err error) {
	queryStr := "SELECT bucket FROM log_sync_settings"
	rows, qErr := glbDB.Query(queryStr)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query failed due to, %s", qErr.Error()))
		return
	}

	buckets = make([]string, 0)
	for rows.Next() {
		var bucket string

		rErr := rows.Scan(&bucket)
		if rErr != nil {
			err = errors.New(fmt.Sprintf("read row data failed due to, %s", rErr.Error()))
			return
		}
		buckets = append(buckets, bucket)
	}
	return
}

//创建日志表
func CreateTableIfNone(bucket string, date string) (err error) {
	stmt, sErr := glbDB.Prepare(GetCreateLogRecordTableSQL(bucket, date))
	if sErr != nil {
		err = errors.New(fmt.Sprintf("prepare exec failed due to, %s", sErr.Error()))
		return
	}
	defer stmt.Close()
	_, execErr := stmt.Exec()
	if execErr != nil {
		err = errors.New(fmt.Sprintf("failed to create table log_record_%s due to, %s", bucket+"_"+date, execErr.Error()))
		return
	}
	return
}

//写入日志记录
func WriteQLogRecord(id string, bucket string, date string, reqIp string, reqTime time.Time, reqMethod string, reqPath string, reqProto string, statusCode int,
	totalBytes int, referer string, userAgent string, host string, version string) (err error) {
	stmt, sErr := glbDB.Prepare(fmt.Sprintf("INSERT INTO %s (id,bucket,date,req_ip,req_time,req_method,req_path,req_proto,status_code,total_bytes,referer,user_agent,host,version) "+
		" VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE id=?", GetLogRecordTableName(bucket, date)))
	if sErr != nil {
		err = errors.New(fmt.Sprintf("prepare exec failed due to, %s", sErr.Error()))
		return
	}
	defer stmt.Close()
	_, execErr := stmt.Exec(id, bucket, date, reqIp, reqTime, reqMethod, reqPath, reqProto, statusCode, totalBytes, referer, userAgent, host, version, id)
	if execErr != nil {
		err = errors.New(fmt.Sprintf("failed to insert or ignore due to, %s", execErr.Error()))
		return
	}
	return
}
