package server

import (
	"database/sql"
	"fmt"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	USER string = com.GetEnvVar("DB_USER")
	PASSWORD string = com.GetEnvVar("DB_PASSWORD")
	ENDPOINT string = com.GetEnvVar("DB_ENDPOINT")
	DATABASE string = com.GetEnvVar("DB_DATABASE")
)

func Start() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", USER, PASSWORD, ENDPOINT, DATABASE)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}

	if err = http.ListenAndServe(com.GetEnvVar("PORT"), handlers(db)); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}

func handlers(db *sql.DB) http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/addresses", addressList(db)).Methods("GET")
	return router
}

