package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-save-water/pkg/log"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type WaterUsage struct {
	AccountNumber int    `json:"accountNumber"`
	BillDate      string `json:"billDate"`
	Usage         string `json:"waterUsage"`
	ImageURL      string `json:"imageURL"`
	CreatedDT     string `json:"createDT"`
	ModifiedDT    string `json:"modifiedDT"`
}

func addBill(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") == "application/json" {

			var newBill WaterUsage
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Error.Println("Error retrieving information")
			} else {
				json.Unmarshal(reqBody, &newBill)
			}

			AccountNumber := newBill.AccountNumber
			BillDate := newBill.BillDate
			ImageURL := ""
			CreatedDT := time.Now().Format(time.RFC3339)
			ModifiedDT := CreatedDT

			query := fmt.Sprintf("INSERT INTO WaterUsage VALUE(%d, '%s', '%s', '%s', '%s')", AccountNumber, BillDate, ImageURL, CreatedDT, ModifiedDT)

			_, err = db.Query(query)
			if err != nil {
				log.Error.Println("Error retrieving information")
			} else {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Bill added successfully"))
			}
		}
	}
}

func getAllBill(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var allUsage []WaterUsage
		results, err := db.Query("SELECT FROM WaterUsage")

		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		for results.Next() {
			var usage WaterUsage
			err = results.Scan(&usage.AccountNumber, &usage.BillDate, &usage.Usage)

			allUsage = append(allUsage, usage)

			if err != nil {
				log.Error.Println("Error retrieving information")
			}

			json.NewEncoder(w).Encode(usage)
		}
	}
}

func getBill(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		BillDate := params["billDate"]
		results, err := db.Query("SELECT FROM WaterUsage WHERE BillDate='%s'", BillDate)

		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		var usage WaterUsage
		err = results.Scan(&usage.AccountNumber, &usage.BillDate, &usage.Usage)

		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		json.NewEncoder(w).Encode(usage)
	}
}

func updateBill(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateBill WaterUsage

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err.Error()
		} else {
			json.Unmarshal(reqBody, &updateBill)
		}

		newBillDate := updateBill.BillDate
		newUsage := updateBill.Usage
		modifiedDT := time.Now().Format(time.RFC3339)

		query := fmt.Sprintf("UPDATE WaterUsage SET BillDate='%s', Usage='%s', ModifiedDT='%s'", newBillDate, newUsage, modifiedDT)
		_, err = db.Query(query)
		if err != nil {
			log.Error.Println("Error updating information")
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("201 - Bill Date: " + newBillDate + "Usage: " + newUsage + "successfully updated."))
	}
}

func deleteBill(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var deleteBill WaterUsage

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err.Error()
		} else {
			json.Unmarshal(reqBody, deleteBill)
		}

		query := fmt.Sprintf("DELETE FROM WaterUsage WHERE BillDate='%s'", deleteBill.BillDate)
		_, err = db.Query(query)
		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("201 - Month Usage: " + deleteBill.BillDate + "successfully deleted"))
	}
}
