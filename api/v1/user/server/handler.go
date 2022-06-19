package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type UsersInfo struct {
	UserID         int    `json:"userID"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashedPassword,omitempty"`
	Email          string `json:"email"`
	Role           string `json:"role,omitempty"`
	IsDeleted      bool   `json:"isDeleted,omitempty"`
	CreatedDT      string `json:"createdDT,omitempty"`
	ModifiedDT     string `json:"modifiedDT,omitempty"`
}

func userList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		results, err := db.Query("SELECT UserID, FirstName, LastName, UserName, Email, Role, IsDeleted, CreatedDT, ModifiedDT  FROM User")
		if err != nil {
			log.Error.Println(err)
		}

		var userList []UsersInfo

		for results.Next() {
			// map this type to the record in the table
			var user UsersInfo
			err = results.Scan(
				&user.UserID,
				&user.FirstName,
				&user.LastName,
				&user.Username,
				&user.Email,
				&user.Role,
				&user.IsDeleted,
				&user.CreatedDT,
				&user.ModifiedDT,
			)
			userList = append(userList, user)
			if err != nil {
				log.Error.Println(err)
			}
		}

		json.NewEncoder(w).Encode(userList)
	}
}

// func courseGet to scan client input courseCode and courseTitle with sql database and return course results
func userGet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user UsersInfo

		params := mux.Vars(r)
		userID := params["userid"]

		// Prepared Statement
		stmt, err := db.Prepare("SELECT UserID, FirstName, LastName, UserName, Email, Role, IsDeleted FROM User WHERE UserID = ?")
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		result := stmt.QueryRow(userID)
		err = result.Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.IsDeleted,
		)

		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No user found"))
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func userPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			newUser       UsersInfo
			userNameExist bool
		)

		// TODO: Use Middleware when merged
		if r.Header.Get("Content-type") == "application/json" {
			// read the string sent to the service
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newUser)

				// TODO: Validation?
				/*if newUser.Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 -Please provide FirstName " + "LastName" + "Username" + "Email"))
					return
				}*/

				// Check if username exist
				err = db.QueryRow("call spUserExistByUserName(?)", newUser.Username).Scan(&userNameExist)
				if err != nil {
					log.Error.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - Internal Server Error"))
					return
				}
				if userNameExist {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Username Taken"))
					return
				}
				//query = fmt.Sprintf("UPDATE user SET username ='%s' WHERE username = '%s'", newUser.first_name, newUser.last_name)
				_, err = db.Query("call spUserCreate(?, ?, ?, ?, ?, ?)",
					newUser.Username,
					newUser.FirstName,
					newUser.LastName,
					newUser.Email,
					newUser.HashedPassword,
					newUser.Role,
				)
				if err != nil {
					log.Error.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - Server Error"))
					return
				}
				w.WriteHeader(http.StatusCreated)
			}
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply data in JSON format"))
		}
	}
}

func userPut(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userID := params["userid"]

		var (
			userInfo  UsersInfo
			userExist bool
		)

		// TODO: Use Middleware when merged
		if r.Header.Get("Content-type") == "application/json" {
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &userInfo)

				// TODO: Validation?
				/*if newUser.Username == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 -Please provide username " + " email " + "in JSON format"))
					return
				}*/

				// Check if userID exist
				err = db.QueryRow("call spUserExistByUserID(?)", userID).Scan(&userExist)
				if err != nil {
					log.Error.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - Internal Server Error"))
					return
				}
				if !userExist {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("404 - User not Found"))
					return
				}

				_, err = db.Query("call spUserUpdate(?, ?, ?, ?, ?, ?, ?, ?)",
					userID,
					com.NewNullString(userInfo.FirstName),
					com.NewNullString(userInfo.LastName),
					com.NewNullString(userInfo.Username),
					com.NewNullString(userInfo.HashedPassword),
					com.NewNullString(userInfo.Email),
					com.NewNullString(userInfo.Role),
					userInfo.IsDeleted,
				)

				if err != nil {
					log.Error.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - Internal Server Error"))
					return
				}

				w.WriteHeader(http.StatusAccepted)

			}
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply user information in JSON format"))
		}
	}
}

// func userDelete to check scan client input Username, FirstName, LastName with the database, delete if matches
func userDelete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userExist bool

		params := mux.Vars(r)
		userID := params["userid"]

		err := db.QueryRow("call spUserExistByUserID(?)", userID).Scan(&userExist)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		if !userExist {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - User not Found"))
			return
		}

		// Prepared Statement
		stmt, err := db.Prepare("UPDATE User SET IsDeleted = true, ModifiedDT = NOW() WHERE UserID = ?")
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		_, err = stmt.Query(userID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 -Server Error"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - User deleted: " + userID))
	}
}
