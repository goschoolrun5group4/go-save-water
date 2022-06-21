package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-save-water/pkg/log"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type WaterUsage struct {
	AccountNumber int    `json:"accountNumber"`
	BillDate      string `json:"billDate"`
	Consumption   string `json:"consumption"`
	ImageURL      string `json:"imageURL,omitempty"`
	CreatedDT     string `json:"createDT,omitempty"`
	ModifiedDT    string `json:"modifiedDT,omitempty"`
}

func addUsage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUsage WaterUsage
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error.Println("Error retrieving information")
		} else {
			json.Unmarshal(reqBody, &newUsage)
		}

		AccountNumber := newUsage.AccountNumber
		BillDate := newUsage.BillDate
		Consumption := newUsage.Consumption
		//ImageURL := ""

		query := fmt.Sprintf("INSERT INTO WaterUsage (AccountNumber, BillDate, Consumption, ImageURL, CreatedDT, ModifiedDT) VALUES(%d, '%s', %s, null, now(), null)", AccountNumber, BillDate, Consumption)

		_, err = db.Query(query)
		if err != nil {
			log.Error.Println("Error retrieving information")
		} else {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Bill added successfully"))
		}
	}
}

func getUsages(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		accountNumber := params["accountNumber"]
		var allUsage []WaterUsage
		query := fmt.Sprintf("SELECT AccountNumber, BillDate, Consumption FROM WaterUsage WHERE AccountNumber=%s", accountNumber)
		results, err := db.Query(query)

		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		for results.Next() {
			var usage WaterUsage
			err = results.Scan(&usage.AccountNumber, &usage.BillDate, &usage.Consumption)

			allUsage = append(allUsage, usage)

			if err != nil {
				log.Error.Println("Error retrieving information")
			}
		}
		json.NewEncoder(w).Encode(allUsage)
	}
}

func getUsage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		accountNumber := params["accountNumber"]
		billDate := params["billDate"]
		query := fmt.Sprintf("SELECT AccountNumber, BillDate, Consumption FROM WaterUsage WHERE AccountNumber=%s AND billDate='%s'", accountNumber, billDate)
		results := db.QueryRow(query)

		var usage WaterUsage
		err := results.Scan(&usage.AccountNumber, &usage.BillDate, &usage.Consumption)

		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		json.NewEncoder(w).Encode(usage)
	}
}

func updateUsage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		accountNumber := params["accountNumber"]
		oldBillDate := params["billDate"]
		var updatedUsage WaterUsage

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err.Error()
		} else {
			json.Unmarshal(reqBody, &updatedUsage)
		}

		newBillDate := updatedUsage.BillDate
		newUsage := updatedUsage.Consumption

		query := fmt.Sprintf("UPDATE WaterUsage SET BillDate='%s', Consumption='%s', ModifiedDT=now() WHERE AccountNumber=%s AND BillDate='%s'", newBillDate, newUsage, accountNumber, oldBillDate)
		_, err = db.Query(query)
		if err != nil {
			log.Error.Println("Error updating information")
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("201 - Bill Date: " + newBillDate + "Usage: " + newUsage + "successfully updated."))
	}
}

func deleteUsage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		accountNumber := params["accountNumber"]
		billDate := params["billDate"]

		query := fmt.Sprintf("DELETE FROM WaterUsage WHERE BillDate='%s' AND AccountNumber=%s", billDate, accountNumber)
		_, err := db.Query(query)
		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("201 - Monthly Usage: " + billDate + "successfully deleted"))
	}
}
