package server

import (
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func Start() {

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(com.GetEnvVar("PORT"), router); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}
