package qlog

import (
	"fmt"
	"strings"
)

func GetLogRecordTableName(bucket string, date string) string {
	return fmt.Sprintf("log_record_%s_%s", bucket, strings.Replace(date, "-", "_", -1))
}

func GetCreateLogRecordTableSQL(bucket string, date string) string {
	tmpl := `CREATE TABLE IF NOT EXISTS [TABLE_NAME] (
  id varchar(60) NOT NULL,
  req_ip varchar(15) DEFAULT NULL,
  req_time datetime DEFAULT NULL,
  req_method varchar(10) DEFAULT NULL,
  req_path varchar(1000) DEFAULT NULL,
  req_proto varchar(15) DEFAULT NULL,
  status_code int(11) DEFAULT NULL,
  total_bytes int(11) DEFAULT NULL,
  referer varchar(2000) DEFAULT NULL,
  user_agent varchar(2000) DEFAULT NULL,
  host varchar(500) DEFAULT NULL,
  version varchar(10) DEFAULT NULL,
  bucket varchar(63) DEFAULT NULL,
  date char(10) DEFAULT NULL,
  ip_code varchar(100) DEFAULT NULL,
  ip_country varchar(100) DEFAULT NULL,
  ip_region varchar(100) DEFAULT NULL,
  ip_city varchar(100) DEFAULT NULL,
  ip_isp varchar(100) DEFAULT NULL,
  ip_note varchar(100) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`
	sql := strings.Replace(tmpl, "[TABLE_NAME]", GetLogRecordTableName(bucket, date), -1)
	return sql
}
