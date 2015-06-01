package qlog

import (
	"net/http"
	"strconv"
)

type RetReport struct {
	FormData map[string]interface{}
	RetData  map[string]interface{}
	Error    string
}

func (this *QLogServer) serveTopAccessResource(w http.ResponseWriter, req *http.Request) {
	templates := []string{
		"views/base.html",
		"views/base_d.html",
		"views/head.html",
		"views/header.html",
		"views/footer.html",
		"views/reports/top_access_resource.html",
	}
	reqMethod := req.Method
	buckets, err := GetBucketListFromSettings()
	if err != nil {
		this.renderHtml(w, &RetReport{Error: err.Error()}, templates)
		return
	}
	formData := make(map[string]interface{}, 0)
	formData["Buckets"] = buckets
	if reqMethod == "GET" {
		this.renderHtml(w, &RetReport{FormData: formData}, templates)
	} else if reqMethod == "POST" {
		bucket := req.FormValue("bucket")
		date := req.FormValue("date")
		num, _ := strconv.Atoi(req.FormValue("num"))
		//check param

		//query
		resResult, err := GetTopAccessResource(bucket, date, num)
		if err != nil {
			this.renderHtml(w, &RetReport{Error: err.Error()}, templates)
		} else {
			retData := make(map[string]interface{}, 0)
			retData["TopAccessResource"] = resResult
			formData["Bucket"] = bucket
			formData["Date"] = date
			formData["Num"] = num
			this.renderHtml(w, &RetReport{FormData: formData, RetData: retData}, templates)
		}
	}
}

func (this *QLogServer) serveSupplierSummary(w http.ResponseWriter, req *http.Request) {
	templates := []string{
		"views/base.html",
		"views/base_d.html",
		"views/head.html",
		"views/header.html",
		"views/footer.html",
		"views/reports/supplier_summary.html",
	}
	reqMethod := req.Method
	buckets, err := GetBucketListFromSettings()
	if err != nil {
		this.renderHtml(w, &RetReport{Error: err.Error()}, templates)
		return
	}
	formData := make(map[string]interface{}, 0)
	formData["Buckets"] = buckets
	if reqMethod == "GET" {
		this.renderHtml(w, &RetReport{FormData: formData}, templates)
	} else if reqMethod == "POST" {
		bucket := req.FormValue("bucket")
		date := req.FormValue("date")
		//query
		accessCntOfSuppliers, err := GetAccessCountOfSuppliers(bucket, date)
		if err != nil {
			this.renderHtml(w, &RetReport{Error: err.Error()}, templates)
		} else {
			retData := make(map[string]interface{}, 0)
			retData["AccessCntOfSuppliers"] = accessCntOfSuppliers
			formData["Bucket"] = bucket
			formData["Date"] = date
			this.renderHtml(w, &RetReport{FormData: formData, RetData: retData}, templates)
		}
	}
}
