package server

import (
	"database/sql"
	"encoding/json"
	"go-save-water/pkg/log"
	"net/http"
)

type RewardInfo struct {
	RewardID    string `json:"rewardID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	RedeemAmt   string `json:"redeemAmt"`
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
