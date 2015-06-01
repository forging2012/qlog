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

//模版函数
var (
	funcMap = template.FuncMap{
		"empty":  Empty,
		"nempty": NotEmpty,
		"urlFor": UrlFor,
	}
)

type QLogServer struct {
	Config *QLogConfig
}

func (this *QLogServer) Listen() error {
	http.HandleFunc("/", this.serveIndex)
	http.HandleFunc("/prepare", this.servePrepare)
	http.HandleFunc("/prepare/query", this.servePrepareQuery)
	http.HandleFunc("/status/query", this.serveStatusQuery)
	http.HandleFunc("/settings", this.serveSettings)
	http.HandleFunc("/settings/add", this.serveSettingsAdd)
	http.HandleFunc("/settings/delete", this.serveSettingsDelete)
	http.HandleFunc("/report/top/access", this.serveTopAccessResource)
	http.HandleFunc("/report/supplier/summary", this.serveSupplierSummary)

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
	Status  []*LogPrepareStatus
	Error   string
}

type RetSettings struct {
	SettingsAll []*QLogSyncSettings
	Error       string
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
	buckets, err := GetBucketListFromSettings()
	if err != nil {
		this.renderHtml(w, RetIndex{Error: err.Error()}, templates)
	} else {
		prepStatus, perr := QueryLogPrepare()
		if perr != nil {
			this.renderHtml(w, RetIndex{Error: perr.Error()}, templates)
		} else {
			this.renderHtml(w, RetIndex{Buckets: buckets, Status: prepStatus}, templates)
		}
	}
}

func (this *QLogServer) serveSettings(w http.ResponseWriter, req *http.Request) {
	templates := []string{
		"views/base.html",
		"views/base_d.html",
		"views/head.html",
		"views/header.html",
		"views/footer.html",
		"views/settings.html",
	}
	settingsAll, err := GetLogSyncSettingsAll()
	if err != nil {
		this.renderHtml(w, RetSettings{Error: err.Error()}, templates)
	} else {
		this.renderHtml(w, RetSettings{SettingsAll: settingsAll}, templates)
	}
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
