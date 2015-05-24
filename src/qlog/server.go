package qlog

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	funcMap = template.FuncMap{
		"empty":  Empty,
		"nempty": NotEmpty,
	}
)

type QLogServer struct {
	Config *QLogConfig
}

func (this *QLogServer) Listen() error {

	http.HandleFunc("/", this.serveIndex)
	http.HandleFunc("/settings", this.serveSettings)
	http.HandleFunc("/settings/add", this.serveSettings)
	http.HandleFunc("/settings/delete", this.serveSettings)
	http.HandleFunc("/settings/edit", this.serveSettings)

	http.HandleFunc("/static/fonts/glyphicons-halflings-regular.woff", this.serveStatic)
	http.HandleFunc("/static/fonts/glyphicons-halflings-regular.wof2", this.serveStatic)
	http.HandleFunc("/static/fonts/glyphicons-halflings-regular.eot", this.serveStatic)
	http.HandleFunc("/static/fonts/glyphicons-halflings-regular.svg", this.serveStatic)
	http.HandleFunc("/static/fonts/glyphicons-halflings-regular.ttf", this.serveStatic)

	http.HandleFunc("/static/css/bootstrap.min.css", this.serveStatic)
	http.HandleFunc("/static/css/main.css", this.serveStatic)

	http.HandleFunc("/static/js/jquery.min.js", this.serveStatic)
	http.HandleFunc("/static/js/bootstrap.min.js", this.serveStatic)
	http.HandleFunc("/static/js/main.js", this.serveStatic)

	endPoint := fmt.Sprintf("%s:%d", this.Config.ListenHost, this.Config.ListenPort)
	server := &http.Server{
		Addr: endPoint,
	}
	listenErr := server.ListenAndServe()
	if listenErr != nil {
		return errors.New(fmt.Sprintf("start server error, %s", listenErr.Error()))
	}
	return nil
}

func (this *QLogServer) serveStatic(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.RequestURI, "/")
	staticFp, openErr := os.Open(path)
	if openErr != nil {
		http.Error(w, openErr.Error(), http.StatusNotFound)
		return
	}
	defer staticFp.Close()
	if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css; cahrset=utf-8")
	}
	_, cpErr := io.Copy(w, staticFp)
	if cpErr != nil {
		http.Error(w, cpErr.Error(), http.StatusInternalServerError)
		return
	}
}

type RetIndex struct {
	Buckets []string
	Records *[]QLogRecord
	Error   string
}

func (this *QLogServer) serveIndex(w http.ResponseWriter, req *http.Request) {
	templates := []string{
		"views/base.html",
		"views/base_d.html",
		"views/head.html",
		"views/header.html",
		"views/footer.html",
		"views/index.html",
	}
	reqMethod := req.Method
	var errMsg string

	buckets, err := GetBucketListFromSettings()
	if err != nil {
		errMsg = err.Error()
	}
	if reqMethod == "GET" {
		this.renderHtml(w, RetIndex{Buckets: buckets, Error: errMsg}, templates)
		return
	}
	if reqMethod == "POST" {
		bucket := req.FormValue("bucket")
		date := req.FormValue("date")
		logType := req.FormValue("type")
		records, err := QueryLogData(bucket, date, []string{logType}...)
		if err != nil {
			this.renderHtml(w, RetIndex{Buckets: buckets, Error: err.Error()}, templates)
			return
		}
		this.renderHtml(w, RetIndex{Buckets: buckets, Records: records, Error: errMsg}, templates)
	}
}

func (this *QLogServer) serveSettings(w http.ResponseWriter, req *http.Request) {

}

func (this *QLogServer) renderHtml(w http.ResponseWriter, data interface{}, viewPaths []string) {
	tmpl := template.New("base.html")
	tmpl = tmpl.Funcs(funcMap)
	tmpl, tErr := tmpl.ParseFiles(viewPaths...)

	if tErr != nil {
		http.Error(w, tErr.Error(), http.StatusInternalServerError)
		return
	}
	eErr := tmpl.Execute(w, data)
	if eErr != nil {
		http.Error(w, eErr.Error(), http.StatusInternalServerError)
	}
}
