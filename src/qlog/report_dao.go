package qlog

import (
	"errors"
	"fmt"
)

func GetTopAccessResource(bucket, date string, num int) (resResult []*TopAccessResource, err error) {
	resResult = make([]*TopAccessResource, 0)
	queryStr := fmt.Sprintf("select host,req_path, concat(host,req_path) as url, count(*) as cnt "+
		"from %s group by req_path order by cnt desc limit 0, %d", GetLogRecordTableName(bucket, date), num)
	rows, qErr := glbDB.Query(queryStr)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query failed due to, %s", qErr.Error()))
		return
	}

	for rows.Next() {
		res := &TopAccessResource{}
		var host string
		var path string
		var url string
		var count int

		rErr := rows.Scan(&host, &path, &url, &count)
		if rErr != nil {
			err = errors.New(fmt.Sprintf("read row data failed due to, %s", rErr.Error()))
			return
		}
		res.Host = host
		res.Path = path
		res.Url = url
		res.Count = count
		resResult = append(resResult, res)
	}
	return
}

func GetAccessCountOfSuppliers(bucket, date string) (resResult []*AccessCntOfSupplier, err error) {
	resResult = make([]*AccessCntOfSupplier, 0)
	queryStr := fmt.Sprintf("select ip_isp,count(ip_isp) as cnt from %s group by ip_isp order by cnt desc",
		GetLogRecordTableName(bucket, date))
	rows, qErr := glbDB.Query(queryStr)
	if qErr != nil {
		err = errors.New(fmt.Sprintf("query failed due to, %s", qErr.Error()))
		return
	}

	for rows.Next() {
		var isp string
		var cnt int

		item := &AccessCntOfSupplier{}
		rErr := rows.Scan(&isp, &cnt)
		if rErr != nil {
			err = errors.New(fmt.Sprintf("read row data failed due to, %s", rErr.Error()))
			return
		}
		item.Supplier = isp
		item.Count = cnt
		if isp == "" {
			item.Supplier = "N/A"
		}
		resResult = append(resResult, item)
	}
	return
}
