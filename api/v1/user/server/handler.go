package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"io/ioutil"
	"net/http"
	"strconv"

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
	IsDeleted      *bool  `json:"isDeleted,omitempty"`
	CreatedDT      string `json:"createdDT,omitempty"`
	ModifiedDT     string `json:"modifiedDT,omitempty"`
	Verified       *bool  `json:"verified,omitempty"`
	PointBalance   int    `json:"pointBalance,omitempty"`
}

type Transaction struct {
	TransactionID string  `json:"transactionID"`
	UserID        string  `json:"userID"`
	Type          string  `json:"type"`
	RewardID      *string `json:"rewardID,omitempty"`
	Points        string  `json:"points"`
	TransactionDT string  `json:"transactionDT"`
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
		stmt, err := db.Prepare("SELECT UserID, FirstName, LastName, UserName, Email, Role, IsDeleted, Verified, PointBalance FROM User WHERE UserID = ?")
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
			&user.Verified,
			&user.PointBalance,
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

func userGetByEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UsersInfo

		params := mux.Vars(r)
		email := params["email"]

		// Prepared Statement
		stmt, err := db.Prepare("SELECT UserID, FirstName, LastName, UserName, Email, Role, IsDeleted, Verified FROM User WHERE Email = ?")
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		result := stmt.QueryRow(email)
		err = result.Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.IsDeleted,
			&user.Verified,
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

			_, err = db.Query("call spUserUpdate(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				userID,
				com.NewNullString(userInfo.FirstName),
				com.NewNullString(userInfo.LastName),
				com.NewNullString(userInfo.Username),
				com.NewNullString(userInfo.HashedPassword),
				com.NewNullString(userInfo.Email),
				com.NewNullString(userInfo.Role),
				userInfo.IsDeleted,
				userInfo.Verified,
				com.NewNullInt64(userInfo.PointBalance),
			)

			if err != nil {
				log.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
				return
			}

			w.WriteHeader(http.StatusAccepted)

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

// func userAddPoints to add point to user table
func userAddPoints(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userID := params["userid"]
		pointsToAddStr := params["points"]
		pointsToAdd, _ := strconv.Atoi(pointsToAddStr)

		// Create a helper function for preparing failure results.
		fail := func(err error, statusCode int, body string) {
			log.Error.Println(err)
			w.WriteHeader(statusCode)
			w.Write([]byte(body))
		}

		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			fail(err, http.StatusInternalServerError, com.MsgServerError)
			return
		}
		// Defer a rollback in case anything fails.
		defer tx.Rollback()

		// Get Current Point Balance
		var pointBalance int

		if err = tx.QueryRowContext(ctx, "SELECT PointBalance FROM User WHERE UserID = ?", userID).Scan(&pointBalance); err != nil {
			if err == sql.ErrNoRows {
				fail(err, http.StatusNotFound, com.MsgUserNotFound)
				return
			}
			fail(err, http.StatusInternalServerError, com.MsgServerError)
			return
		}

		newPointBalance := pointBalance + pointsToAdd

		// Update Point
		_, err = tx.ExecContext(ctx, "UPDATE User SET PointBalance = ?, ModifiedDT = NOW() WHERE UserID = ?", newPointBalance, userID)
		if err != nil {
			fail(err, http.StatusInternalServerError, com.MsgServerError)
			return
		}

		// Create Entry in Transaction
		_, err = tx.ExecContext(ctx, "INSERT INTO Transaction (UserID, Type, Points) VALUES (?, ?, ?)", userID, "Earn", pointsToAdd)
		if err != nil {
			fail(err, http.StatusInternalServerError, com.MsgServerError)
			return
		}

		// Commit the transaction.
		if err = tx.Commit(); err != nil {
			fail(err, http.StatusInternalServerError, com.MsgServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func userTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		userID := params["userid"]

		// Create a helper function for preparing failure results.
		fail := func(err error, statusCode int, body string) {
			log.Error.Println(err)
			w.WriteHeader(statusCode)
			w.Write([]byte(body))
		}

		var transactions []Transaction
		stmt, err := db.Prepare("SELECT TransactionID, Type, RewardID, Points, TransactionDT FROM Transaction WHERE UserID = ? ORDER BY TransactionDT DESC")
		if err != nil {
			fail(err, http.StatusInternalServerError, com.MsgServerError)
			return
		}

		results, err := stmt.Query(userID)
		if err != nil {
			fail(err, http.StatusInternalServerError, com.MsgServerError)
			return
		}

		for results.Next() {
			var transaction Transaction
			err = results.Scan(&transaction.TransactionID, &transaction.Type, &transaction.RewardID, &transaction.Points, &transaction.TransactionDT)
			if err != nil {
				fail(err, http.StatusInternalServerError, com.MsgServerError)
				return
			}
			transactions = append(transactions, transaction)
		}

		json.NewEncoder(w).Encode(transactions)
	}
}
