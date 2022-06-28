package server

import (
	"database/sql"
	"net/http"

	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	mw "go-save-water/pkg/middleware"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func Start() {

	db := dbConnect()
	defer db.Close()

	log.Info.Println("Server Start")

	if err := http.ListenAndServe(com.GetEnvVar("PORT"), handlers(db)); err != nil {
		log.Fatal.Fatalln("ListenAndServe: ", err)
	}
}

func handlers(db *sql.DB) http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	std := alice.New(mw.ContentTypeHandler)

	api.Handle("/users", userList(db))
	api.Handle("/user", std.Then(userPost(db))).Methods("POST")
	api.Handle("/user/{userid:[0-9]+}", userGet(db)).Methods("GET")
	api.Handle("/user/email/{email}", userGetByEmail(db)).Methods("GET")
	api.Handle("/user/{userid:[0-9]+}", std.Then(userPut(db))).Methods("PUT")
	api.Handle("/user/{userid:[0-9]+}", userDelete(db)).Methods("DELETE")
	api.Handle("/user/{userid:[0-9]+}", userDelete(db)).Methods("DELETE")
	api.Handle("/user/{userid:[0-9]+}/points/{points:[0-9]+}", userAddPoints(db)).Methods("POST")

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
