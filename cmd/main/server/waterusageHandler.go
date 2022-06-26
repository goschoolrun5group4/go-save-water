package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"io/ioutil"
	"net/http"
)

func getUsages(w http.ResponseWriter, r *http.Request) {
	type Results struct {
		AccountNumber int
		BillDate      string
		Consumption   string
	}

	viewData := struct {
		Error       bool
		ErrorMsg    string
		ClientMsg   string
		ShowResults []Results
	}{
		false,
		"",
		"",
		nil,
	}

	if r.Method == http.MethodPost {
		accountNumber := r.FormValue("accountNumber")

		url := com.GetEnvVar("API_USAGE_ADDR") + fmt.Sprintf("/usages/%s", accountNumber)
		req, err := http.NewRequest("GET", url, nil)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Error.Println(err)
		}

		if err != nil {
			viewData.Error = true
			viewData.ErrorMsg = "Internal server error"
		}

		if res.StatusCode == http.StatusUnauthorized {
			viewData.Error = true
			viewData.ErrorMsg = "Unable to delete current selection"
		}

		if res.StatusCode == http.StatusUnprocessableEntity {
			viewData.Error = true
			viewData.ErrorMsg = "Please enter valid date"
		}

		if res.StatusCode == http.StatusNotFound {
			viewData.Error = true
			viewData.ErrorMsg = "Data not found"
		}

		if !viewData.Error {
			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				viewData.Error = true
				viewData.ErrorMsg = "Internal server error"
			} else {
				json.Unmarshal(body, &viewData.ShowResults)
				fmt.Println(viewData.ShowResults)

				if err != nil {
					viewData.Error = true
					viewData.ErrorMsg = "Internal server error"
				} else {
					res.StatusCode = http.StatusAccepted
					viewData.ClientMsg = fmt.Sprintf("User Found")
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "getusages.gohtml", viewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func getUsage(w http.ResponseWriter, r *http.Request) {
	viewData := struct {
		Error     bool
		ErrorMsg  string
		ClientMsg string
		Usage     map[string]interface{}
	}{
		false,
		"",
		"",
		nil,
	}

	if r.Method == http.MethodPost {
		accountNumber := r.FormValue("accountNumber")
		searchDate := r.FormValue("searchDate")

		url := com.GetEnvVar("API_USAGE_ADDR") + fmt.Sprintf("/usage/%s/%s", accountNumber, searchDate)
		//jsonValue := fmt.Sprintf(`{"accountNumber":%s,"billDate":"%s"}`, accountNumber, searchDate)

		req, err := http.NewRequest("GET", url, nil)
		//req.Header.Set("Content-type", "application/json") sending json data

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Error.Println(err)
		}

		if res.StatusCode == http.StatusUnauthorized {
			viewData.Error = true
			viewData.ErrorMsg = "Unable to delete current selection"
		}

		if res.StatusCode == http.StatusUnprocessableEntity {
			viewData.Error = true
			viewData.ErrorMsg = "Please enter valid date"
		}

		if res.StatusCode == http.StatusNotFound {
			viewData.Error = true
			viewData.ErrorMsg = "Data not found"
		}

		if !viewData.Error {
			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				viewData.Error = true
				viewData.ErrorMsg = "Internal server error"
			} else {
				json.Unmarshal(body, &viewData.Usage)

				if err != nil {
					viewData.Error = true
					viewData.ErrorMsg = "Internal server error"
				} else {
					res.StatusCode = http.StatusAccepted
					viewData.Error = false
					viewData.ClientMsg = fmt.Sprintf("Date Found")
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "getusage.gohtml", viewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func addUsages(w http.ResponseWriter, r *http.Request) {
	viewData := struct {
		Error     bool
		ErrorMsg  string
		ClientMsg string
	}{
		false,
		"",
		"",
	}

	if r.Method == http.MethodPost {
		accountNumber := r.FormValue("accountNumber")
		newDate := r.FormValue("newDate")
		newConsumption := r.FormValue("newConsumption")
		newImageURL := r.FormValue("newImageURL")

		url := com.GetEnvVar("API_USAGE_ADDR") + "/usage"
		jsonValue := fmt.Sprintf(`{"accountNumber":%s,"billDate":"%s","consumption":"%s","imageURL":"%s"}`, accountNumber, newDate, newConsumption, newImageURL)

		var jsonData = []byte(jsonValue)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Error.Println(err)
		}

		if res.StatusCode == http.StatusUnauthorized {
			viewData.Error = true
			viewData.ErrorMsg = "Unable to add new usages"
		}

		if res.StatusCode == http.StatusConflict {
			viewData.Error = true
			viewData.ErrorMsg = "Date already exist"
		}

		if res.StatusCode == http.StatusUnprocessableEntity {
			viewData.Error = true
			viewData.ErrorMsg = "Unable to leave fields blank"
		}

		if !viewData.Error {
			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				viewData.Error = true
				viewData.ErrorMsg = "Internal server error"
			} else {
				var addNewBill map[string]interface{}
				json.Unmarshal(body, &addNewBill)

				//modidifiedDT := updatedWaterUsage["modifiedDT"]

				if err != nil {
					viewData.Error = true
					viewData.ErrorMsg = "Internal server error"
				} else {
					res.StatusCode = http.StatusAccepted
					viewData.ClientMsg = fmt.Sprintf("Date: %s, Usage: %s, bill image successfully added", newDate, newConsumption)
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "addusage.gohtml", viewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func updateUsage(w http.ResponseWriter, r *http.Request) {
	viewData := struct {
		Error     bool
		ErrorMsg  string
		ClientMsg string
	}{
		false,
		"",
		"",
	}

	if r.Method == http.MethodPost {
		accountNumber := r.FormValue("accountNumber")
		oldDate := r.FormValue("oldDate")
		updateDate := r.FormValue("updateDate")
		updateConsumption := r.FormValue("updateConsumption")
		updateImageURL := r.FormValue("updateImageURL")

		url := com.GetEnvVar("API_USAGE_ADDR") + fmt.Sprintf("/usage/%s/%s", accountNumber, oldDate)
		jsonVal := fmt.Sprintf(`{"billDate":"%s","consumption":"%s","imageURL":"%s"}`, updateDate, updateConsumption, updateImageURL)

		var jsonData = []byte(jsonVal)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Error.Println(err)
		}

		if res.StatusCode == http.StatusConflict {
			viewData.Error = true
			viewData.ErrorMsg = "Date already exist"
		}

		if res.StatusCode == http.StatusUnauthorized {
			viewData.Error = true
			viewData.ErrorMsg = "Unable to delete current selection"
		}

		if res.StatusCode == http.StatusUnprocessableEntity {
			viewData.Error = true
			viewData.ErrorMsg = "Please enter valid date"
		}

		if res.StatusCode == http.StatusNotFound {
			viewData.Error = true
			viewData.ErrorMsg = "Data not found"
		}

		if !viewData.Error {
			_, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				viewData.Error = true
				viewData.ErrorMsg = "Please enter valid date"
			} else {
				if err != nil {
					viewData.Error = true
					viewData.ErrorMsg = "Internal server error"
				} else {
					res.StatusCode = http.StatusAccepted
					viewData.ClientMsg = fmt.Sprintf("New Date: %s , New Usage: %s and new bill successfully updated", updateDate, updateConsumption)
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "updateusage.gohtml", viewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func deleteUsage(w http.ResponseWriter, r *http.Request) {
	viewData := struct {
		Error     bool
		ErrorMsg  string
		ClientMsg string
		Usage     map[string]interface{}
	}{
		false,
		"",
		"",
		nil,
	}

	if r.Method == http.MethodPost {
		accountNumber := r.FormValue("accountNumber")
		deleteDate := r.FormValue("deleteDate")

		url := com.GetEnvVar("API_USAGE_ADDR") + fmt.Sprintf("/usage/%s/%s", accountNumber, deleteDate)
		//jsonValue := fmt.Sprintln(`{"accountNumber":%s,"billDate":"%s"}`, accountNumber, deleteDate)

		//var jsonData = []byte(jsonValue)
		req, err := http.NewRequest("DELETE", url, nil)
		//req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Error.Println(err)
		}

		if res.StatusCode == http.StatusUnauthorized {
			viewData.Error = true
			viewData.ErrorMsg = "Unable to delete current selection"
		}

		if res.StatusCode == http.StatusUnprocessableEntity {
			viewData.Error = true
			viewData.ErrorMsg = "Please enter valid date"
		}

		if res.StatusCode == http.StatusInternalServerError {
			viewData.Error = true
			viewData.ErrorMsg = "Internal server error"
		}

		if res.StatusCode == http.StatusNotFound {
			viewData.Error = true
			viewData.ErrorMsg = "Data not found"
		}

		if !viewData.Error {
			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				viewData.Error = true
				viewData.ErrorMsg = "Internal server error"
			} else {
				var deletedWaterUsage map[string]interface{}
				json.Unmarshal(body, &deletedWaterUsage)

				//modidifiedDT := updatedWaterUsage["modifiedDT"]

				if err != nil {
					viewData.Error = true
					viewData.ErrorMsg = "Internal server error"
				} else {
					res.StatusCode = http.StatusAccepted
					viewData.ClientMsg = fmt.Sprintf("Usage Date: %s successfully deleted", deleteDate)
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "deleteusage.gohtml", viewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}
