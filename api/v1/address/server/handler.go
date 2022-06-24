package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type AddressInfo struct {
	AccountNumber int    `json:"accountNumber"`
	PostalCode    string `json:"postalCode"`
	Floor         string `json:"floor"`
	UnitNumber    string `json:"unitNumber"`
	BuildingName  string `json:"buildingName"`
	BlockNumber   string `json:"blockNumber"`
	Street        string `json:"street"`
	CreatedDT     string `json:"createdDT"`
	ModifiedDT    string `json:"modifiedDT"`
}

func createAddress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var newAddress AddressInfo
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &newAddress)

			PostalCode := newAddress.PostalCode
			if len(PostalCode) != 6 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Minimum length for PostalCode is 6"))
				return
			}

			Floor := newAddress.Floor
			UnitNumber := newAddress.UnitNumber
			BuildingName := newAddress.BuildingName
			BlockNumber := newAddress.BlockNumber
			Street := newAddress.Street
			CreatedDt := time.Now().Format(time.RFC3339)
			ModifiedDt := CreatedDt

			query := fmt.Sprintf(
				"INSERT INTO Address VALUES (NULL, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
				PostalCode, Floor, UnitNumber, BuildingName, BlockNumber, CreatedDt, ModifiedDt, Street)
			_, err := db.Query(query)
			if err != nil {
				panic(err.Error())
			}

			lastId := fmt.Sprintf("SELECT LAST_INSERT_ID()")

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Address Account Number: " + lastId + " added successfully"))
		}
	}
}

func updateAddress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var newAddress AddressInfo
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &newAddress)
			AccountNumber := newAddress.AccountNumber

			if AccountNumber < 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please select the address you want to update"))
				return
			}

			PostalCode := newAddress.PostalCode
			if len(PostalCode) != 6 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Minimum length for PostalCode is 6"))
				return
			}

			Floor := newAddress.Floor
			UnitNumber := newAddress.UnitNumber
			BuildingName := newAddress.BuildingName
			BlockNumber := newAddress.BlockNumber
			Street := newAddress.Street
			ModifiedDt := time.Now().Format(time.RFC3339)

			query := fmt.Sprintf(
				"UPDATE Address SET PostalCode='%s', Floor='%s', UnitNumber='%s', BuildingName='%s', BlockNumber='%s', ModifiedDT='%s', Street=='%s', WHERE AccountNumber=%d",
				PostalCode, Floor, UnitNumber, BuildingName, BlockNumber, ModifiedDt, Street,  AccountNumber)
			_, err := db.Query(query)
			if err != nil {
				panic(err.Error())
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Address Account Number: " + strconv.Itoa(AccountNumber) + " updated successfully"))
		}
	}
}

func readAddresses(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := db.Query("Select * FROM Address")

		if err != nil {
			panic(err.Error())
		}

		var addresses []AddressInfo

		for results.Next() {
			var address AddressInfo
			err = results.Scan(
				&address.AccountNumber,
				&address.PostalCode,
				&address.Floor,
				&address.UnitNumber,
				&address.BuildingName,
				&address.BlockNumber,
				&address.Street,
				&address.CreatedDT,
				&address.ModifiedDT,
			)

			addresses = append(addresses, address)

			if err != nil {
				panic(err.Error())
			}

			json.NewEncoder(w).Encode(addresses)
		}
	}
}

func readAddress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		AccountNumber := params["accountnumber"]
		results, err := db.Query("Select * FROM Address WHERE AccountNumber=%d", AccountNumber)

		if err != nil {
			panic(err.Error())
		}

		var address AddressInfo
		err = results.Scan(
			&address.AccountNumber,
			&address.PostalCode,
			&address.Floor,
			&address.UnitNumber,
			&address.BuildingName,
			&address.CreatedDT,
			&address.ModifiedDT,
			&address.Street,
		)

		if err != nil {
			panic(err.Error())
		}

		json.NewEncoder(w).Encode(address)
	}
}

func deleteAddress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var address AddressInfo
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &address)
			AccountNumber := address.AccountNumber

			if AccountNumber < 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please select the address you want to delete"))
				return
			}

			query := fmt.Sprintf("DELETE FROM Address WHERE AccountNumber='%s'", AccountNumber)
			_, err := db.Query(query)
			if err != nil {
				panic(err.Error())
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - AccountNumber: " + strconv.Itoa(AccountNumber) + " deleted successfully"))
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please select the address you want to delete"))
		}
	}
}
