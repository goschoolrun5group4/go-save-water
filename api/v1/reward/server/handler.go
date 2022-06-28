package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type RewardInfo struct {
	RewardID    string `json:"rewardID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	RedeemAmt   int    `json:"redeemAmt"`
}

type RedeemInfo struct {
	UserID   int `json:"userID"`
	RewardID int `json:"rewardID"`
}

type ResponseInfo struct {
	StatusCode int
	Body       string
}

func rewards(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rewards []RewardInfo
		results, err := db.Query("SELECT RewardID, Title, Description, Quantity, RedeemAmt FROM Reward WHERE IsDeleted = false")
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		for results.Next() {
			var reward RewardInfo
			err = results.Scan(&reward.RewardID, &reward.Title, &reward.Description, &reward.Quantity, &reward.RedeemAmt)
			if err != nil {
				log.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
				return
			}
			rewards = append(rewards, reward)
		}

		if len(rewards) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not found"))
			return
		}

		json.NewEncoder(w).Encode(rewards)
	}
}

func reward(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		rewardID := params["rewardID"]
		var reward RewardInfo
		stmt, err := db.Prepare("SELECT RewardID, Title, Description, Quantity, RedeemAmt FROM Reward WHERE IsDeleted = false AND RewardID = ?")
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		result := stmt.QueryRow(rewardID)
		err = result.Scan(&reward.RewardID, &reward.Title, &reward.Description, &reward.Quantity, &reward.RedeemAmt)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - No user found"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		json.NewEncoder(w).Encode(reward)
	}
}

func redeem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		var redeemInfo RedeemInfo
		json.Unmarshal(reqBody, &redeemInfo)

		chn := make(chan ResponseInfo)
		go processRedemption(db, redeemInfo, chn)
		resp := <-chn
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(resp.Body))
		return
	}
}

func processRedemption(db *sql.DB, redeemInfo RedeemInfo, chn chan ResponseInfo) {
	//time.Sleep(10 * time.Second)

	// Check if Award is fully redeem
	var (
		reward   RewardInfo
		respInfo ResponseInfo
		userInfo map[string]interface{}
	)
	// Get User Info
	url := com.GetEnvVar("API_USER_ADDR") + fmt.Sprintf("/user/%d", redeemInfo.UserID)
	body, statusCode, err := com.FetchData(url)
	if err != nil {
		log.Error.Println(err)
		if statusCode == http.StatusNotFound {
			respInfo.StatusCode = http.StatusNotFound
			respInfo.Body = "404 - User not found"
			chn <- respInfo
			return
		}
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}
	json.Unmarshal(body, &userInfo)

	stmt, err := db.Prepare("SELECT Quantity, RedeemAmt FROM Reward WHERE IsDeleted = false AND RewardID = ?")
	if err != nil {
		log.Error.Println(err)
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}
	result := stmt.QueryRow(redeemInfo.RewardID)
	err = result.Scan(&reward.Quantity, &reward.RedeemAmt)
	if err != nil {
		log.Error.Println(err)
		if err == sql.ErrNoRows {
			respInfo.StatusCode = http.StatusNotFound
			respInfo.Body = "404 - Award not found"
			chn <- respInfo
			return
		}
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}

	// If user have 0 points
	if userInfo["pointBalance"] == nil {
		respInfo.StatusCode = http.StatusBadRequest
		respInfo.Body = "400 - Not enough point for redemption"
		chn <- respInfo
		return
	}

	userPointBalance := int(userInfo["pointBalance"].(float64))
	// If user don't have enough points to redeem
	if userPointBalance < reward.RedeemAmt {
		respInfo.StatusCode = http.StatusBadRequest
		respInfo.Body = "400 - Not enough point for redemption"
		chn <- respInfo
		return
	}

	// If reward is fully redeem return status conflict
	if reward.Quantity == 0 {
		respInfo.StatusCode = http.StatusConflict
		respInfo.Body = "409 - Fully Redeem"
		chn <- respInfo
		return
	}

	// Proceed with redemption
	stmt, err = db.Prepare("UPDATE Reward SET Quantity = ? WHERE RewardID = ?")
	if err != nil {
		log.Error.Println(err)
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}

	_, err = stmt.Query(reward.Quantity-1, redeemInfo.RewardID)
	if err != nil {
		log.Error.Println(err)
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}

	// Update User's Point
	newUserPointBalance := userPointBalance - reward.RedeemAmt
	jsonStr := fmt.Sprintf("{\"pointBalance\":%d}", newUserPointBalance)

	url = com.GetEnvVar("API_USER_ADDR") + fmt.Sprintf("/user/%d", redeemInfo.UserID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Error.Println(err)
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}

	// Create Entry in Transaction
	stmt, err = db.Prepare("INSERT INTO Transaction (UserID, Type, RewardID, Points) VALUE (?, ?, ?, ?)")
	if err != nil {
		log.Error.Println(err)
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}

	_, err = stmt.Query(redeemInfo.UserID, "Redemption", redeemInfo.RewardID, reward.RedeemAmt)
	if err != nil {
		log.Error.Println(err)
		respInfo.StatusCode = http.StatusInternalServerError
		respInfo.Body = "500 - Internal Server Error"
		chn <- respInfo
		return
	}

	respInfo.StatusCode = http.StatusAccepted
	respInfo.Body = "202 - Accepted"
	chn <- respInfo
}
