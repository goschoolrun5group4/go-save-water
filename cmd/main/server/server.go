package server

import (
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	tpl *template.Template
	fm  = template.FuncMap{
		"formatDateTime": com.FormatDateTime,
	}
)

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*"))
}

func Start() {

	log.Info.Println("Server Start")

	router := mux.NewRouter()

	router.HandleFunc("/", index)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/setup", setup)
	router.HandleFunc("/address/edit", addressEdit)
	router.HandleFunc("/user/edit", userEdit)
	router.HandleFunc("/dashboard", dashboard)
	//router.HandleFunc("/usages", getUsages)
	router.HandleFunc("/usage", usage)
	//router.HandleFunc("/getusage", getUsage)
	//router.HandleFunc("/addusage", addUsages)
	//router.HandleFunc("/updateusage", updateUsage)
	//router.HandleFunc("/deleteusage", deleteUsage)
	router.HandleFunc("/verification/{token}", verification)
	router.HandleFunc("/reward/{rewardID:[0-9]+}", rewardDetail)
	router.HandleFunc("/transaction", transactions)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	if err := http.ListenAndServe(com.GetEnvVar("PORT"), router); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}
