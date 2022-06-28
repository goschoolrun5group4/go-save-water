package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	com "go-save-water/pkg/common"
	"io/ioutil"
	"net/http"

	"go-save-water/pkg/log"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserID         int    `json:"userID,omitempty"`
	Username       string `json:"username"`
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Password       string `json:"password,omitempty"`
	HashedPassword string `json:"hashedPassword,omitempty"`
	Email          string `json:"email,omitempty"`
	Role           string `json:"role,omitempty"`
	SessionID      string `json:"sessionID,omitempty"`
	ExpireDT       string `json:"expireDT,omitempty"`
	AccountNumber  string `json:"accountNumber,omitempty"`
	PointBalance   string `json:"pointBalance"`
}

func signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var signupUser UserInfo

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &signupUser)
		signupUser.Role = "user"

		bPassword, err := bcrypt.GenerateFromPassword([]byte(signupUser.Password), bcrypt.MinCost)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}
		signupUser.HashedPassword = string(bPassword)
		signupUser.Password = ""

		// Call User Create API
		jsonStr, err := json.Marshal(signupUser)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		url := com.GetEnvVar("API_USER_ADDR") + "/user"
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
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("409 - Username Taken"))
			return
		}

		if resp.StatusCode == http.StatusUnprocessableEntity {
			log.Error.Println("SignUp not sending json to user create service.")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
		}

		if resp.StatusCode == http.StatusInternalServerError {
			log.Error.Println("Error Processing from User create service")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
		}

	}
}

func verification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			verificationJson map[string]interface{}
			userInfo         map[string]interface{}
			loginUser        UserInfo
		)

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &verificationJson)

		// Check if user Exist
		url := com.GetEnvVar("API_USER_ADDR") + "/user/email/" + verificationJson["email"].(string)
		body, statusCode, err := com.FetchData(url)
		if err != nil {
			log.Error.Println(err)
			if statusCode == http.StatusNotFound {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Not found"))
				return
			}
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		json.Unmarshal(body, &userInfo)

		if !userInfo["verified"].(bool) {
			// Update User
			jsonStr := "{\"verified\":true, \"pointBalance\": 1000}"

			url = fmt.Sprintf("%s/user/%.0f", com.GetEnvVar("API_USER_ADDR"), userInfo["userID"])
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(jsonStr)))
			req.Header.Set("Content-type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
				return
			}

			if resp.StatusCode != http.StatusAccepted {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
				return
			}
		}

		// Create Entry in Transaction
		stmt, err := db.Prepare("INSERT INTO Transaction (UserID, Type, Points) VALUE (?, ?, ?)")
		_, err = stmt.Query(userInfo["userID"], "Earn", 1000)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		// Process Session Data, delete old session if multiple sessions are detected.
		result := db.QueryRow("CALL spUserSessionCreate(?)", userInfo["username"].(string))
		err = result.Scan(&loginUser.UserID, &loginUser.Username, &loginUser.SessionID, &loginUser.ExpireDT)
		// If user don't exist
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

		json.NewEncoder(w).Encode(loginUser)
	}
}

func login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			loginUser      UserInfo
			hashedPassword string
			verified       bool
			email          string
		)

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &loginUser)

		// Get user password from DB
		result := db.QueryRow("CALL spAuthenticationGet(?)", loginUser.Username)
		err = result.Scan(&hashedPassword, &verified, &email)
		// If user don't exist
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect username or password."))
			return
		}

		// If user is not verified
		if !verified {
			w.WriteHeader(http.StatusBadRequest)
			ret := fmt.Sprintf("{\"email\":\"%s\"}", email)
			w.Write([]byte(ret))
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginUser.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect username or password."))
			return
		}
		loginUser.Password = ""

		// Process Session Data, delete old session if multiple sessions are detected.
		result = db.QueryRow("CALL spUserSessionCreate(?)", loginUser.Username)
		err = result.Scan(&loginUser.UserID, &loginUser.Username, &loginUser.SessionID, &loginUser.ExpireDT)
		// If user don't exist
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

		json.NewEncoder(w).Encode(loginUser)
	}
}

func logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type LogoutUser struct {
			UserID    int    `json:"userID"`
			SessionID string `json:"sessionID"`
		}

		var logoutUser LogoutUser

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &logoutUser)

		// Delete user session
		_, err = db.Query("call spUserSessionDelete(?, ?)", logoutUser.UserID, logoutUser.SessionID)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

	}
}

func verifySession(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		sessionID := params["sessionID"]

		var (
			userInfo   UserInfo
			accountNum sql.NullString
		)
		results := db.QueryRow("CALL spUserSessionGet(?)", sessionID)
		err := results.Scan(&userInfo.UserID, &userInfo.Username, &userInfo.FirstName, &userInfo.LastName, &userInfo.Email, &userInfo.Role, &accountNum, &userInfo.PointBalance)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not found"))
			return
		}

		if accountNum.Valid {
			userInfo.AccountNumber = accountNum.String
		}

		json.NewEncoder(w).Encode(userInfo)
	}
}
