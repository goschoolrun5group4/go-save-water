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

	router.HandleFunc("/", index)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/login", login)
	router.HandleFunc("/dashboard", dashboard)
	router.HandleFunc("/usages", getUsages)
	router.HandleFunc("/usage", getUsage)
	router.HandleFunc("/getusage", getUsage)
	router.HandleFunc("/addusage", addUsages)
	router.HandleFunc("/updateusage", updateUsage)
	router.HandleFunc("/deleteusage", deleteUsage)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	if err := http.ListenAndServe(com.GetEnvVar("PORT"), router); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}
