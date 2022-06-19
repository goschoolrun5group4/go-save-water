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

	api.HandleFunc("/users", allusers)
	api.HandleFunc("/user/{userid}", userGet).Methods("GET")
	api.HandleFunc("/user/{userid}", userPut).Methods("PUT")
	api.HandleFunc("/user/{userid}", userPost).Methods("POST")
	api.HandleFunc("/user/{userid}", userDelete).Methods("DELETE")
	return router
}
