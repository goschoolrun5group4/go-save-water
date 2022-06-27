package server

import (
	"database/sql"
	"fmt"
	"net/http"

	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	mw "go-save-water/pkg/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var (
	USER     string = com.GetEnvVar("DB_USER")
	PASSWORD string = com.GetEnvVar("DB_PASSWORD")
	ENDPOINT string = com.GetEnvVar("DB_ENDPOINT")
	DATABASE string = com.GetEnvVar("DB_DATABASE")
)

func Start() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", USER, PASSWORD, ENDPOINT, DATABASE)
	db, err := sql.Open("mysql", connectionString)

	defer db.Close()

	if err != nil {
		panic(err.Error())
	}

	log.Info.Println("Server Start")

	if err = http.ListenAndServe(com.GetEnvVar("PORT"), handlers(db)); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}

func handlers(db *sql.DB) http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	std := alice.New(mw.ContentTypeHandler)

	api.Handle("/address", std.Then(createAddress(db))).Methods("POST")
	api.Handle("/addresses", readAddresses(db)).Methods("GET")
	api.Handle("/address/{accountnumber}", readAddress(db)).Methods("GET")
	api.Handle("/address/{accountnumber}", std.Then(updateAddress(db))).Methods("PUT")
	api.Handle("/address/{accountnumber}", deleteAddress(db)).Methods("DELETE")

	return router
}
