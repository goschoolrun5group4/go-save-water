package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type UsersInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

func allusers(w http.ResponseWriter, r *http.Request) {
	var user UsersInfo

	results, err := db.Query("Select * FROM user")
	if err != nil {
		panic(err.Error())
	}

	var userList []UsersInfo

	for results.Next() {
		// map this type to the record in the table
		var user UsersInfo
		err = results.Scan(&user.FirstName, &user.LastName, &user.Username, &user.Email)
		userList = append(userList, user)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(user)
}

// func courseGet to scan client input courseCode and courseTitle with sql database and return course results
func userGet(w http.ResponseWriter, r *http.Request) {

	var user UsersInfo

	params := mux.Vars(r)
	userID := params["userid"]

	query := fmt.Sprintf("Select * FROM user where email = '%s' ", userID)

	result := db.QueryRow(query)

	err := result.Scan(&user.FirstName, &user.LastName, &user.Username, &user.Email)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 -No user found"))
		return
	}
	json.NewEncoder(w).Encode(user)

}

func userPost(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userID := params["userid"]

	if r.Header.Get("Content-type") == "application/json" {
		// read the string sent to the service
		var newUser UsersInfo
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &newUser)
			if newUser.Username == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 -Please provide FirstName " + "LastName" + "Username" + "Email"))
				return
			}
			var user UsersInfo

			query := fmt.Sprintf("Select * FROM Username where userID = '%s' ", userID)

			result := db.QueryRow(query)

			err := result.Scan(&user.FirstName, &user.LastName, &user.Username, &user.Email)

			if err == nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("409 -Duplicate user info"))
				return
			}
			//query = fmt.Sprintf("UPDATE user SET username ='%s' WHERE username = '%s'", newUser.first_name, newUser.last_name)
			query = fmt.Sprintf("INSERT INTO user VALUES ('%s', '%s')", userID, newUser.first_name, newUser.last_name, newUser.username, newUser.Email)
			_, err = db.Query(query)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 -Server Error"))
				return
			}
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 -Please supply username " + "in JSON format"))
		}
	}
}

func userPut(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userID := params["userid"]

	var newUser UsersInfo
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		json.Unmarshal(reqBody, &newUser)
		if newUser.Username == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 -Please provide username " + " email " + "in JSON format"))
			return
		}
		var user UsersInfo

		query := fmt.Sprintf("Select * FROM username where email = '%s' ", userID)

		result := db.QueryRow(query)

		err := result.Scan(&course.CourseCode, &course.CourseTitle)

		if err == nil {
			query = fmt.Sprintf("UPDATE user SET email  = '%s' WHERE username = '%s'", newUser.UsersInfo, userID)
			_, err = db.Query(query)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 -Server Error"))
				return
			}
		} else {
			query = fmt.Sprintf("INSERT INTO user VALUES ('%s', '%s')", userID, newUser.UsersInfo)
			_, err = db.Query(query)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 -Server Error"))
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 -Please supply " + "user information " + "in JSON format"))
	}
}

// func userDelete to check scan client input Username, FirstName, LastName with the database, delete if matches
func userDelete(w http.ResponseWriter, r *http.Request) {

	var user UsersInfo

	params := mux.Vars(r)
	userID := params["userid"]

	query := fmt.Sprintf("Select * FROM user where username = '%s' ", userID)

	result := db.QueryRow(query)

	err := result.Scan(&user.Username, &user.Email)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 -No course found"))
		return
	}
	query = fmt.Sprintf("DELETE FROM username WHERE email = '%s'", userID)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 -Server Error"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("202 -User deleted: " + params["userid"]))
}
