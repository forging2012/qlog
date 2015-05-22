package qlog

import (
	"errors"
	"fmt"
	"html/template"

	"net/http"
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

func (this *QLogServer) serveIndex(w http.ResponseWriter, req *http.Request) {
	templates := []string{
		"views/base.html",
		"views/base_d.html",
		"views/head.html",
		"views/header.html",
		"views/footer.html",
		"views/index.html",
	}
	this.renderHtml(w, "hello world", templates)
}

func (this *QLogServer) serveSettings(w http.ResponseWriter, req *http.Request) {

}

func (this *QLogServer) renderHtml(w http.ResponseWriter, data interface{}, viewPaths []string) {
	tmpl, tErr := template.ParseFiles(viewPaths...)
	if tErr != nil {
		http.Error(w, tErr.Error(), http.StatusInternalServerError)
		return
	}
	eErr := tmpl.Execute(w, data)
	if eErr != nil {
		http.Error(w, eErr.Error(), http.StatusInternalServerError)
	}
}
