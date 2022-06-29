package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-save-water/pkg/log"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type WaterUsage struct {
	AccountNumber int    `json:"accountNumber,omitempty"`
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
			log.Error.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not found"))
			return
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
		query := fmt.Sprintf("SELECT AccountNumber, BillDate, Consumption FROM WaterUsage WHERE AccountNumber=%s ORDER BY BillDate DESC", accountNumber)
		results, err := db.Query(query)

		if err != nil {
			log.Error.Println("Error retrieving information")
		}

		for results.Next() {
			var usage WaterUsage
			err = results.Scan(&usage.AccountNumber, &usage.BillDate, &usage.Consumption)

			allUsage = append(allUsage, usage)

			if err != nil {
				log.Error.Println(err)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Not found"))
				return
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
			log.Error.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not found"))
			return
		}

		json.NewEncoder(w).Encode(usage)
	}
}

func updateUsage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		accountNumber := params["accountNumber"]
		oldBillDate := params["billDate"]

		query := fmt.Sprintf("SELECT COUNT(*) FROM WaterUsage WHERE AccountNumber=%s AND BillDate='%s'", accountNumber, oldBillDate)
		results := db.QueryRow(query)

		var count int

		err := results.Scan(&count)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		if count == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not found"))
			return
		}

		var updatedUsage WaterUsage

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err.Error()
		} else {
			json.Unmarshal(reqBody, &updatedUsage)
		}

		newBillDate := updatedUsage.BillDate
		newUsage := updatedUsage.Consumption

		query = fmt.Sprintf("UPDATE WaterUsage SET BillDate='%s', Consumption='%s', ModifiedDT=now() WHERE AccountNumber=%s AND BillDate='%s'", newBillDate, newUsage, accountNumber, oldBillDate)
		_, err = db.Query(query)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("201 - Update Successful"))
	}
}

func deleteUsage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		accountNumber := params["accountNumber"]
		billDate := params["billDate"]

		query := fmt.Sprintf("SELECT COUNT(*) FROM WaterUsage WHERE AccountNumber=%s AND BillDate='%s'", accountNumber, billDate)
		results := db.QueryRow(query)

		var count int

		err := results.Scan(&count)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		if count == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not found"))
			return
		}

		query = fmt.Sprintf("DELETE FROM WaterUsage WHERE BillDate='%s' AND AccountNumber=%s", billDate, accountNumber)
		_, err = db.Query(query)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("201 - Monthly Usage: " + billDate + "successfully deleted"))
	}
}

func getUsageByLatestMonths(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		accountNumber := params["accountNumber"]
		numberOfMonths, _ := strconv.Atoi(params["numOfMths"])

		startDate, endDate := getLatestMonths(numberOfMonths)

		results, err := db.Query("CALL spWaterUsageGetByDateRange(?, ?, ?)", accountNumber, startDate, endDate)
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Server Error"))
			return
		}

		var usageList []WaterUsage

		for results.Next() {
			// map this type to the record in the table
			var usage WaterUsage
			err = results.Scan(
				&usage.BillDate,
				&usage.Consumption,
			)
			if err != nil {
				log.Error.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Server Error"))
				return
			}
			usageList = append(usageList, usage)
		}

		json.NewEncoder(w).Encode(usageList)
	}
}

func getNationalUsageByLatestMonths(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		numberOfMonths, _ := strconv.Atoi(params["numOfMths"])

		chn := make(chan *WaterUsage)

		var usageList []*WaterUsage
		for i := 0; i < numberOfMonths; i++ {
			go getNationalAveData(db, i, chn)
		}

		for i := 0; i < numberOfMonths; i++ {
			usage := <-chn
			if usage != nil {
				usageList = append(usageList, usage)
			}
		}

		sort.Slice(usageList, func(i, j int) bool {
			return usageList[i].BillDate < usageList[j].BillDate
		})

		json.NewEncoder(w).Encode(usageList)
	}
}

func getLatestMonths(numOfMths int) (time.Time, time.Time) {

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	endDate := firstOfMonth.AddDate(0, 1, -1)

	start := time.Now()
	startDate := start.AddDate(0, -numOfMths+1, 0)
	startYear, startMonth, _ := startDate.Date()
	startDate = time.Date(startYear, startMonth, 1, 0, 0, 0, 0, time.UTC)

	return startDate, endDate
}

func getNationalAveData(db *sql.DB, numOfMths int, chn chan *WaterUsage) {
	now := time.Now()

	newDate := now.AddDate(0, -numOfMths, -4)
	startYear, startMonth, _ := newDate.Date()
	startDate := time.Date(startYear, startMonth, 1, 0, 0, 0, 0, time.UTC)

	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	endDate := firstOfMonth.AddDate(0, 1-numOfMths, -1)

	results := db.QueryRow("CALL spNationalWaterUsageGetByDateRange(?, ?)", startDate, endDate)
	var (
		sqlDate        sql.NullString
		sqlConsumption sql.NullString
	)

	err := results.Scan(&sqlDate, &sqlConsumption)

	if err != nil {
		log.Error.Println(err)
		chn <- nil
	}

	if sqlDate.Valid && sqlConsumption.Valid {
		usage := &WaterUsage{
			BillDate:    sqlDate.String,
			Consumption: sqlConsumption.String,
		}
		chn <- usage
	} else {
		chn <- nil
	}
}
