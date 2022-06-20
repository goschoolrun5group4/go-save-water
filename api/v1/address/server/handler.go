package server

import (
	"database/sql"
	"fmt"
	"time"
)

func addressList(db *sql.DB) {
	// ReadAddresses(db *sql.DB)
}

type Address struct {
	AccountNumber int
	PostalCode int
	Floor int
	UnitNumber int
	BuildingName string
	BlockNumber string
	CreatedDT string
	ModifiedDT string
}

func CreateAddress(db *sql.DB, PostalCode int, Floor int, UnitNumber int, BuildingName string) {
	CreatedDt := time.Now().Format(time.RFC3339)
	ModifiedDt := CreatedDt

	query := fmt.Sprintf(
		"INSERT INTO Address VALUES (%d, %d, %d, '%s', '%s', '%s')", 
		PostalCode, Floor, UnitNumber, BuildingName, CreatedDt, ModifiedDt)    
	_, err := db.Query(query)  
	if err != nil {
		panic(err.Error())
	}
}

func ReadAddresses(db *sql.DB) {   
	results, err := db.Query("Select * FROM Address")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var address Address
		err = results.Scan(
			&address.AccountNumber, 
			&address.PostalCode, 
			&address.Floor, 
			&address.UnitNumber,
			&address.BuildingName,
			&address.CreatedDT,
			&address.ModifiedDT)

		if err != nil {
			panic(err.Error()) 
		}      
		
		fmt.Println(
			address.AccountNumber, 
			address.PostalCode, 
			address.Floor, 
			address.UnitNumber,
			address.BuildingName,
			address.CreatedDT,
			address.ModifiedDT)
	}
}

func UpdateAddress(db *sql.DB, AccountNumber int, PostalCode int, Floor int, UnitNumber int, BuildingName string) {
	ModifiedDt := time.Now().Format(time.RFC3339)

	query := fmt.Sprintf(
		"UPDATE Address SET PostalCode=%d, Floor=%d, UnitNumber=%d BuildingName='%s' ModifiedDT='%s' WHERE AccountNumber=%d", 
		PostalCode, Floor, UnitNumber, BuildingName, ModifiedDt, AccountNumber)
	_, err := db.Query(query)   
	if err != nil {
		panic(err.Error())
	}
}

func DeleteAddress(db *sql.DB, AccountNumber int) {
	query := fmt.Sprintf("DELETE FROM Address WHERE AccountNumber='%d'", AccountNumber)
	_, err := db.Query(query)   
	if err != nil {
		panic(err.Error())
	}
}