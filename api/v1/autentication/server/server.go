package server

import (
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	if err := http.ListenAndServe(com.GetEnvVar("PORT"), handlers()); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}

func handlers() http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/login", login).Methods("GET")

	return router
}
