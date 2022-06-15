package server

import (
	"go-save-water/pkg/log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatal.Fatalln(err)
	}
}
