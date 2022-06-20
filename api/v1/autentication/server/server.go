package server

import (
	"database/sql"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	mw "go-save-water/pkg/middleware"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
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

	std := alice.New(mw.ContentTypeHandler)

	api.Handle("/signup", std.Then(signup(db))).Methods("POST")
	api.Handle("/login", std.Then(login(db))).Methods("POST")
	api.Handle("/logout", std.Then(logout(db))).Methods("POST")

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
