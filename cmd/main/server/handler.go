package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"

	"go-save-water/pkg/validator"
)

// createNewSecureCookie creates and return a new secure cookie.
func createNewSecureCookie(uuid string, expireDT time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     com.GetEnvVar("COOKIE_NAME"),
		Expires:  expireDT,
		Value:    uuid,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
	}
	return cookie
}

func index(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {

	type SignupUser struct {
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password"`
		Email     string `json:"email"`
	}

	ViewData := struct {
		ComparePasswordFail bool
		SignupUser          SignupUser
		Error               bool
		UsernameTaken       bool
		ValidateUserName    bool
		ValidateFirstName   bool
		ValidateLastName    bool
		ValidateEmail       bool
		ValidatePassword    bool
		ValidateFail        bool
	}{
		false,
		SignupUser{},
		false,
		false,
		true,
		true,
		true,
		true,
		true,
		false,
	}

	if r.Method == http.MethodPost {
		ViewData.SignupUser.FirstName = r.FormValue("firstName")
		ViewData.SignupUser.LastName = r.FormValue("lastName")
		ViewData.SignupUser.Username = r.FormValue("userName")
		ViewData.SignupUser.Email = r.FormValue("emailAddr")
		ViewData.SignupUser.Password = r.FormValue("password")
		passwordConfirm := r.FormValue("confirmPassword")

		// Validation
		// Validate username
		if validator.IsEmpty(ViewData.SignupUser.Username) || !validator.IsValidUsername(ViewData.SignupUser.Username) {
			ViewData.ValidateUserName = false
			ViewData.ValidateFail = true
		}
		// Validate first name
		if validator.IsEmpty(ViewData.SignupUser.FirstName) || !validator.IsValidName(ViewData.SignupUser.FirstName) {
			ViewData.ValidateFirstName = false
			ViewData.ValidateFail = true
		}
		// Validate last name
		if validator.IsEmpty(ViewData.SignupUser.LastName) || !validator.IsValidName(ViewData.SignupUser.LastName) {
			ViewData.ValidateLastName = false
			ViewData.ValidateFail = true
		}
		// Validate email
		if validator.IsEmpty(ViewData.SignupUser.Email) || !validator.IsValidEmail(ViewData.SignupUser.Email) {
			ViewData.ValidateEmail = false
			ViewData.ValidateFail = true
		}
		// Validate password
		if validator.IsEmpty(ViewData.SignupUser.Password) || !validator.IsValidPassword(ViewData.SignupUser.Password) {
			ViewData.ValidatePassword = false
			ViewData.ValidateFail = true
		}

		// Compare if password is the same
		if c := strings.Compare(ViewData.SignupUser.Password, passwordConfirm); c != 0 {
			ViewData.ComparePasswordFail = true
			ViewData.ValidateFail = true
		}

		if ViewData.ValidateFail {
			if err := tpl.ExecuteTemplate(w, "signup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		jsonStr, err := json.Marshal(ViewData.SignupUser)
		if err != nil {
			log.Error.Println(err)
			return
		}

		url := com.GetEnvVar("API_AUTHENTICATION_ADDR") + "/signup"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		if resp.StatusCode == http.StatusConflict {
			ViewData.ValidateUserName = true
			ViewData.UsernameTaken = true
		}

		if resp.StatusCode == http.StatusInternalServerError {
			ViewData.Error = true
		}

		if !ViewData.Error && !ViewData.ValidateUserName {
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				ViewData.Error = true
			} else {
				var loginInfo map[string]interface{}
				json.Unmarshal(body, &loginInfo)

				uuid := loginInfo["sessionID"].(string)
				date := loginInfo["expireDT"].(string)
				expireDT, err := time.Parse(time.RFC3339, date)

				if err != nil {
					ViewData.Error = true
				} else {
					cookie := createNewSecureCookie(uuid, expireDT)
					http.SetCookie(w, cookie)
					http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				}
			}
		}

	}

	if err := tpl.ExecuteTemplate(w, "signup.gohtml", ViewData); err != nil {
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

		if resp.StatusCode == http.StatusUnauthorized {
			ViewData.Error = true
			ViewData.ErrorMsg = "Incorrect username or password."
		}

		if resp.StatusCode == http.StatusInternalServerError {
			ViewData.Error = true
			ViewData.ErrorMsg = "Internal Server Error."
		}

		if !ViewData.Error {
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				ViewData.Error = true
				ViewData.ErrorMsg = "Internal Server Error."
			} else {
				var loginInfo map[string]interface{}
				json.Unmarshal(body, &loginInfo)

				uuid := loginInfo["sessionID"].(string)
				date := loginInfo["expireDT"].(string)
				expireDT, err := time.Parse(time.RFC3339, date)

				if err != nil {
					ViewData.Error = true
					ViewData.ErrorMsg = "Internal Server Error."
				} else {
					cookie := createNewSecureCookie(uuid, expireDT)
					http.SetCookie(w, cookie)
					http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				}
			}
		}
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
