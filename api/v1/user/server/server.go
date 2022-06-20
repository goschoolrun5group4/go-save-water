package server

import (
	"database/sql"
	"net/http"

	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func Start() {

	db := dbConnect()

	if err := http.ListenAndServe(com.GetEnvVar("PORT"), handlers(db)); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}

func handlers(db *sql.DB) http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	api.Handle("/users", userList(db))
	api.Handle("/user", userPost(db)).Methods("POST")
	api.Handle("/user/{userid}", userGet(db)).Methods("GET")
	api.Handle("/user/{userid}", userPut(db)).Methods("PUT")
	api.Handle("/user/{userid}", userDelete(db)).Methods("DELETE")

	return router
}

func dbConnect() *sql.DB {
	dbCfg := mysql.Config{
		User:                 com.GetEnvVar("DBUSER"),
		Passwd:               com.GetEnvVar("DBPASS"),
		Net:                  "tcp",
		Addr:                 com.GetEnvVar("DBHOST") + ":" + com.GetEnvVar("DBPORT"),
		DBName:               com.GetEnvVar("DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", dbCfg.FormatDSN())

	if err != nil {
		log.Fatal.Fatalln(err)
	}

	return db
}
