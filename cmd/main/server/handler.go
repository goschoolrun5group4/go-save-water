package server

import (
	"bytes"
	"fmt"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	ViewData := struct {
		Error    bool
		ErrorMsg string
	}{
		false,
		"",
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		url := com.GetEnvVar("API_AUTHENTICATION_ADDR") + "/login"
		jsonVal := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)

		var jsonStr = []byte(jsonVal)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error.Println(err)
		}

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			ViewData.Error = true
			ViewData.ErrorMsg = "Incorrect username or password."
		case http.StatusInternalServerError:
			ViewData.Error = true
			ViewData.ErrorMsg = "Internal Server Error."
		default:
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		defer resp.Body.Close()
	}

	if err := tpl.ExecuteTemplate(w, "login.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "dashboard.gohtml", nil); err != nil {
		log.Fatal.Fatalln(err)
	}
}
