package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"go-save-water/pkg/validator"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {

	type SignupUser struct {
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password"`
		Email     string `json:"email"`
	}

	ViewData := struct {
		SignupUser            SignupUser
		Error                 bool
		UsernameTaken         bool
		ErrValidateUserName   bool
		ErrValidateFirstName  bool
		ErrValidateLastName   bool
		ErrValidateEmail      bool
		ErrValidatePassword   bool
		ComparePasswordFail   bool
		ValidateFail          bool
		VerificationEmailSent bool
	}{
		SignupUser{},
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
	}

	if r.Method == http.MethodPost {
		ViewData.SignupUser.FirstName = r.FormValue("firstName")
		ViewData.SignupUser.LastName = r.FormValue("lastName")
		ViewData.SignupUser.Username = r.FormValue("userName")
		ViewData.SignupUser.Email = r.FormValue("emailAddr")
		ViewData.SignupUser.Password = r.FormValue("password")
		passwordConfirm := r.FormValue("confirmPassword")

		// Validation
		// Validate username
		if validator.IsEmpty(ViewData.SignupUser.Username) || !validator.IsValidUsername(ViewData.SignupUser.Username) {
			ViewData.ErrValidateUserName = true
			ViewData.ValidateFail = true
		}
		// Validate first name
		if validator.IsEmpty(ViewData.SignupUser.FirstName) || !validator.IsValidName(ViewData.SignupUser.FirstName) {
			ViewData.ErrValidateFirstName = true
			ViewData.ValidateFail = true
		}
		// Validate last name
		if validator.IsEmpty(ViewData.SignupUser.LastName) || !validator.IsValidName(ViewData.SignupUser.LastName) {
			ViewData.ErrValidateLastName = true
			ViewData.ValidateFail = true
		}
		// Validate email
		if validator.IsEmpty(ViewData.SignupUser.Email) || !validator.IsValidEmail(ViewData.SignupUser.Email) {
			ViewData.ErrValidateEmail = true
			ViewData.ValidateFail = true
		}
		// Validate password
		if validator.IsEmpty(ViewData.SignupUser.Password) || !validator.IsValidPassword(ViewData.SignupUser.Password) {
			ViewData.ErrValidatePassword = true
			ViewData.ValidateFail = true
		}

		// Compare if password is the same
		if c := strings.Compare(ViewData.SignupUser.Password, passwordConfirm); c != 0 {
			ViewData.ComparePasswordFail = true
			ViewData.ValidateFail = true
		}

		if ViewData.ValidateFail {
			if err := tpl.ExecuteTemplate(w, "signup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		jsonStr, err := json.Marshal(ViewData.SignupUser)
		if err != nil {
			log.Error.Println(err)
			return
		}

		url := com.GetEnvVar("API_AUTHENTICATION_ADDR") + "/signup"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Error.Println(err)
			ViewData.Error = true
		}

		if resp.StatusCode == http.StatusConflict {
			ViewData.ErrValidateUserName = true
			ViewData.UsernameTaken = true
		}

		if resp.StatusCode == http.StatusInternalServerError {
			ViewData.Error = true
		}

		if !ViewData.Error && !ViewData.ErrValidateUserName {
			// Go Routine
			//Send email
			go sendVerificationEmail(ViewData.SignupUser.Email)
			ViewData.VerificationEmailSent = true
		}

	}

	if err := tpl.ExecuteTemplate(w, "signup.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func verification(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tokenString := params["token"]

	isValid, email, err := validateJWT(tokenString)

	ViewData := struct {
		Email        string
		TokenExpired bool
		EmailSend    bool
	}{
		email,
		false,
		false,
	}

	if isValid {

		jsonStr := fmt.Sprintf("{\"email\":\"%s\"}", email)

		url := com.GetEnvVar("API_AUTHENTICATION_ADDR") + "/verification"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		// If User not found
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
			if err := tpl.ExecuteTemplate(w, "verification.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Error.Println(err)
		}

		var loginInfo map[string]interface{}
		json.Unmarshal(body, &loginInfo)

		uuid := loginInfo["sessionID"].(string)
		date := loginInfo["expireDT"].(string)
		expireDT, err := time.Parse(time.RFC3339, date)

		if err != nil {
			if err := tpl.ExecuteTemplate(w, "verification.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		} else {
			cookie := createNewSecureCookie(uuid, expireDT)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}
	}

	errCode := err.(*jwt.ValidationError).Errors
	if errCode == jwt.ValidationErrorExpired {
		ViewData.TokenExpired = true
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		fmt.Println(email)
		// Go Routine to send email
		go sendVerificationEmail(email)
		ViewData.EmailSend = true
	}

	if err := tpl.ExecuteTemplate(w, "verification.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	ViewData := struct {
		Error                 bool
		ErrorMsg              string
		ErrUserNotVerified    bool
		Email                 string
		VerificationEmailSent bool
	}{
		false,
		"",
		false,
		"",
		false,
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		url := com.GetEnvVar("API_AUTHENTICATION_ADDR") + "/login"
		jsonVal := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonVal)))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Error.Println(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			ViewData.Error = true
			ViewData.ErrorMsg = "Internal Server Error."
		}

		if resp.StatusCode == http.StatusBadRequest {
			var data map[string]interface{}
			json.Unmarshal(body, &data)
			ViewData.ErrUserNotVerified = true
			go sendVerificationEmail(data["email"].(string))
		}

		if resp.StatusCode == http.StatusUnauthorized {
			ViewData.Error = true
			ViewData.ErrorMsg = "Incorrect username or password."
		}

		if resp.StatusCode == http.StatusInternalServerError {
			ViewData.Error = true
			ViewData.ErrorMsg = "Internal Server Error."
		}

		if !ViewData.Error && !ViewData.ErrUserNotVerified {
			var loginInfo map[string]interface{}
			json.Unmarshal(body, &loginInfo)

			uuid := loginInfo["sessionID"].(string)
			date := loginInfo["expireDT"].(string)
			expireDT, err := time.Parse(time.RFC3339, date)

			if err != nil {
				ViewData.Error = true
				ViewData.ErrorMsg = "Internal Server Error."
			} else {
				cookie := createNewSecureCookie(uuid, expireDT)
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "login.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic.Println(err)
		}
	}()

	loggedInUser, err := authenticationCheck(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	cookie, err := r.Cookie(com.GetEnvVar("COOKIE_NAME"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	url := com.GetEnvVar("API_AUTHENTICATION_ADDR") + "/logout"
	jsonVal := fmt.Sprintf(`{"userID":%d,"sessionID":"%s"}`, int(loggedInUser["userID"].(float64)), cookie.Value)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonVal)))
	req.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Error.Println(err)
	}

	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic.Println(err)
		}
	}()

	loggedInUser, err := authenticationCheck(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser  map[string]interface{}
		UserUsage     string
		NationalUsage string
		UpdateAddress bool
		Rewards       []map[string]interface{}
	}{
		loggedInUser,
		"",
		"",
		false,
		nil,
	}

	if loggedInUser["accountNumber"] == nil {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	chnUserUsage := make(chan string)
	chnNationalUsage := make(chan string)
	chnRewards := make(chan []map[string]interface{})

	go getUserUsage(loggedInUser["accountNumber"].(string), chnUserUsage)
	go getNationalUsage(chnNationalUsage)
	go getRewards(chnRewards)

	for i := 0; i < 3; i++ {
		select {
		case userUsageJson := <-chnUserUsage:
			if len(userUsageJson) > 0 {
				ViewData.UserUsage = userUsageJson
			}
		case nationalUsageJson := <-chnNationalUsage:
			if len(nationalUsageJson) > 0 {
				ViewData.NationalUsage = nationalUsageJson
			}
		case rewards := <-chnRewards:
			ViewData.Rewards = rewards
		}
	}

	if err := tpl.ExecuteTemplate(w, "dashboard.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func setup(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic.Println(err)
		}
	}()

	loggedInUser, err := authenticationCheck(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	type AddressInfo struct {
		UserID        int    `json:"userID"`
		AccountNumber int    `json:"accountNumber"`
		PostalCode    string `json:"postalCode"`
		Floor         string `json:"floor"`
		UnitNumber    string `json:"unitNumber"`
		BuildingName  string `json:"buildingName,omitempty"`
		BlockNumber   string `json:"blockNumber"`
		Street        string `json:"street"`
	}

	type WaterUsage struct {
		AccountNumber int    `json:"accountNumber"`
		BillDate      string `json:"billDate"`
		Consumption   string `json:"consumption"`
	}

	var (
		addressInfo AddressInfo
		usage       WaterUsage
	)

	ViewData := struct {
		Error                   bool
		ValidateConsumptionFail bool
	}{
		false,
		false,
	}

	if r.Method == http.MethodPost {
		accountNum, err := strconv.Atoi(r.FormValue("accountNum"))
		if err != nil {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		// Address
		addressInfo.UserID = int(loggedInUser["userID"].(float64))
		addressInfo.AccountNumber = accountNum
		addressInfo.PostalCode = r.FormValue("postalCode")
		addressInfo.Floor = r.FormValue("floor")
		addressInfo.UnitNumber = r.FormValue("unitNumber")
		addressInfo.BuildingName = r.FormValue("buildingName")
		addressInfo.BlockNumber = r.FormValue("blockNumber")
		addressInfo.Street = r.FormValue("street")
		// Usage
		usage.AccountNumber = accountNum
		usage.BillDate = r.FormValue("billDate")
		usage.Consumption = r.FormValue("consumption")

		// Check if Consumption is valid
		if _, err := strconv.ParseFloat(r.FormValue("consumption"), 64); err != nil {
			ViewData.ValidateConsumptionFail = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		jsonStr, err := json.Marshal(addressInfo)
		if err != nil {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		url := com.GetEnvVar("API_ADDRESS_ADDR") + "/address"
		resp, err := postToApi(url, jsonStr)

		if err != nil {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		if resp.StatusCode == http.StatusInternalServerError {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		jsonStr, err = json.Marshal(usage)
		if err != nil {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		url = com.GetEnvVar("API_USAGE_ADDR") + "/usage"
		resp, err = postToApi(url, jsonStr)
		if err != nil {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		if resp.StatusCode == http.StatusUnauthorized {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		if resp.StatusCode == http.StatusConflict {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		if resp.StatusCode == http.StatusUnprocessableEntity {
			ViewData.Error = true
			if err := tpl.ExecuteTemplate(w, "setup.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return

	}

	if err := tpl.ExecuteTemplate(w, "setup.gohtml", nil); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func addressEdit(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic.Println(err)
		}
	}()

	loggedInUser, err := authenticationCheck(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	type AddressInfo struct {
		PostalCode   string `json:"postalCode"`
		Floor        string `json:"floor"`
		UnitNumber   string `json:"unitNumber"`
		BuildingName string `json:"buildingName"`
		BlockNumber  string `json:"blockNumber"`
		Street       string `json:"street"`
	}

	ViewData := struct {
		LoggedInUser      map[string]interface{}
		AddressInfo       AddressInfo
		ValidateFail      bool
		RetrieveDataError bool
		ProcessError      bool
		ProcessSuccess    bool
	}{
		loggedInUser,
		AddressInfo{},
		false,
		false,
		false,
		false,
	}

	url := com.GetEnvVar("API_ADDRESS_ADDR") + fmt.Sprintf("/address/%s", loggedInUser["accountNumber"].(string))
	body, _, err := com.FetchData(url)
	if err != nil {
		log.Error.Println(err)
		ViewData.RetrieveDataError = true
		if err := tpl.ExecuteTemplate(w, "addressEdit.gohtml", ViewData); err != nil {
			log.Fatal.Fatalln(err)
		}
		return
	}
	err = json.Unmarshal(body, &ViewData.AddressInfo)
	if err != nil {
		log.Error.Println(err)
	}

	if r.Method == http.MethodPost {
		ViewData.AddressInfo.PostalCode = r.FormValue("postalCode")
		ViewData.AddressInfo.Floor = r.FormValue("floor")
		ViewData.AddressInfo.UnitNumber = r.FormValue("unitNumber")
		ViewData.AddressInfo.BuildingName = r.FormValue("buildingName")
		ViewData.AddressInfo.BlockNumber = r.FormValue("blockNumber")
		ViewData.AddressInfo.Street = r.FormValue("street")

		jsonStr, err := json.Marshal(ViewData.AddressInfo)
		if err != nil {
			log.Error.Println(err)
			return
		}

		url := com.GetEnvVar("API_ADDRESS_ADDR") + fmt.Sprintf("/address/%s", loggedInUser["accountNumber"].(string))
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		defer res.Body.Close()
		if err != nil {
			ViewData.ProcessError = true
			if err := tpl.ExecuteTemplate(w, "addressEdit.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		if res.StatusCode == http.StatusInternalServerError {
			ViewData.ProcessError = true
		}
		if res.StatusCode == http.StatusAccepted {
			ViewData.ProcessSuccess = true
		}

	}

	if err := tpl.ExecuteTemplate(w, "addressEdit.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func userEdit(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic.Println(err)
		}
	}()

	type EditUser struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"hashedPassword,omitempty"`
		Email     string `json:"email"`
	}

	loggedInUser, err := authenticationCheck(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var editUser EditUser
	editUser.FirstName = loggedInUser["firstName"].(string)
	editUser.LastName = loggedInUser["lastName"].(string)
	editUser.Email = loggedInUser["email"].(string)

	ViewData := struct {
		LoggedInUser         map[string]interface{}
		EditUserData         EditUser
		ErrValidateFirstName bool
		ErrValidateLastName  bool
		ErrValidateEmail     bool
		ErrValidatePassword  bool
		ComparePasswordFail  bool
		ValidateFail         bool
		ProcessError         bool
		ProcessSuccess       bool
	}{
		loggedInUser,
		editUser,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
	}

	if r.Method == http.MethodPost {
		ViewData.EditUserData.FirstName = r.FormValue("firstName")
		ViewData.EditUserData.LastName = r.FormValue("lastName")
		ViewData.EditUserData.Email = r.FormValue("emailAddr")
		ViewData.EditUserData.Password = r.FormValue("password")
		passwordConfirm := r.FormValue("confirmPassword")

		// Validation
		// Validate first name
		if validator.IsEmpty(ViewData.EditUserData.FirstName) || !validator.IsValidName(ViewData.EditUserData.FirstName) {
			ViewData.ErrValidateFirstName = true
			ViewData.ValidateFail = true
		}
		// Validate last name
		if validator.IsEmpty(ViewData.EditUserData.LastName) || !validator.IsValidName(ViewData.EditUserData.LastName) {
			ViewData.ErrValidateLastName = true
			ViewData.ValidateFail = true
		}
		// Validate email
		if validator.IsEmpty(ViewData.EditUserData.Email) || !validator.IsValidEmail(ViewData.EditUserData.Email) {
			ViewData.ErrValidateEmail = true
			ViewData.ValidateFail = true
		}
		// Validate password
		if len(ViewData.EditUserData.Password) > 0 {
			if !validator.IsValidPassword(ViewData.EditUserData.Password) {
				ViewData.ErrValidatePassword = true
				ViewData.ValidateFail = true
			}
		}

		// Compare if password is the same
		if c := strings.Compare(ViewData.EditUserData.Password, passwordConfirm); c != 0 {
			ViewData.ComparePasswordFail = true
			ViewData.ValidateFail = true
		}

		if ViewData.ValidateFail {
			if err := tpl.ExecuteTemplate(w, "userEdit.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		if len(ViewData.EditUserData.Password) > 0 {
			bPassword, err := bcrypt.GenerateFromPassword([]byte(ViewData.EditUserData.Password), bcrypt.MinCost)
			if err != nil {
				log.Error.Println(err)
				ViewData.ProcessError = true
				if err := tpl.ExecuteTemplate(w, "userEdit.gohtml", ViewData); err != nil {
					log.Fatal.Fatalln(err)
				}
				return
			}
			ViewData.EditUserData.Password = string(bPassword)
		}

		jsonStr, err := json.Marshal(ViewData.EditUserData)
		if err != nil {
			log.Error.Println(err)
			return
		}

		url := com.GetEnvVar("API_USER_ADDR") + fmt.Sprintf("/user/%d", int(loggedInUser["userID"].(float64)))
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Error.Println(err)
			ViewData.ProcessError = true
			if err := tpl.ExecuteTemplate(w, "userEdit.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		switch resp.StatusCode {
		case http.StatusNotFound:
			ViewData.ProcessError = true
		case http.StatusInternalServerError:
			ViewData.ProcessError = true
		default:
			ViewData.ProcessSuccess = true
			loggedInUser, err = authenticationCheck(w, r)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			ViewData.LoggedInUser = loggedInUser
		}

	}

	if err := tpl.ExecuteTemplate(w, "userEdit.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func rewardDetail(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic.Println(err)
		}
	}()

	loggedInUser, err := authenticationCheck(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(r)
	rewardID := params["rewardID"]

	ViewData := struct {
		LoggedInUser            map[string]interface{}
		Rewards                 []map[string]interface{}
		Reward                  map[string]interface{}
		RewardProcessError      bool
		RewardProcessErrorMsg   string
		RewardProcessSuccessful bool
	}{
		loggedInUser,
		nil,
		nil,
		false,
		"",
		false,
	}

	chnReward := make(chan map[string]interface{})
	chnRewards := make(chan []map[string]interface{})

	go getReward(rewardID, chnReward)
	go getRewards(chnRewards)

	for i := 0; i < 2; i++ {
		select {
		case reward := <-chnReward:
			ViewData.Reward = reward
		case rewards := <-chnRewards:
			if rewards != nil {
				rewards = removeCurrReward(rewards, rewardID)
			}
			ViewData.Rewards = rewards
		}
	}

	if r.Method == http.MethodPost {
		checkingErr := false
		pointBalance, err := strconv.Atoi(loggedInUser["pointBalance"].(string))
		if err != nil {
			checkingErr = true
			ViewData.RewardProcessError = true
			ViewData.RewardProcessErrorMsg = "Error processing redemption, please try again later."
		}

		if pointBalance < int(ViewData.Reward["redeemAmt"].(float64)) {
			checkingErr = true
			ViewData.RewardProcessError = true
			ViewData.RewardProcessErrorMsg = "Not enough points to redeem this reward."
		}

		if checkingErr {
			if err := tpl.ExecuteTemplate(w, "rewardDetail.gohtml", ViewData); err != nil {
				log.Fatal.Fatalln(err)
			}
			return
		}

		url := com.GetEnvVar("API_REWARD_ADDR") + "/reward/redeem"
		jsonVal := fmt.Sprintf(`{"userID":%0.f,"rewardID":%s}`, loggedInUser["userID"].(float64), ViewData.Reward["rewardID"].(string))
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonVal)))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Error.Println(err)
		}

		switch resp.StatusCode {
		case http.StatusBadRequest:
			ViewData.RewardProcessError = true
			ViewData.RewardProcessErrorMsg = "Not enough points to redeem this voucher."
		case http.StatusConflict:
			ViewData.RewardProcessError = true
			ViewData.RewardProcessErrorMsg = "Reward has been fully redeemed."
		case http.StatusInternalServerError:
			ViewData.RewardProcessError = true
			ViewData.RewardProcessErrorMsg = "Error Processing the redemption, please try again."
		case http.StatusNotFound:
			ViewData.RewardProcessError = true
			ViewData.RewardProcessErrorMsg = "Error Processing the redemption, please try again."
		default:
			ViewData.RewardProcessSuccessful = true
			// Refresh Data
			chnReward = make(chan map[string]interface{})
			go getReward(rewardID, chnReward)
			reward := <-chnReward
			ViewData.Reward = reward
			// Get User Data
			loggedInUser, err = authenticationCheck(w, r)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			ViewData.LoggedInUser = loggedInUser

		}
	}

	if err := tpl.ExecuteTemplate(w, "rewardDetail.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}

func transactions(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic.Println(err)
		}
	}()

	loggedInUser, err := authenticationCheck(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser map[string]interface{}
		Transactions []map[string]interface{}
	}{
		loggedInUser,
		nil,
	}

	url := com.GetEnvVar("API_USER_ADDR") + fmt.Sprintf("/user/%d/transactions", int(loggedInUser["userID"].(float64)))
	body, _, err := com.FetchData(url)
	if err != nil {
		log.Error.Println(err)
		if err := tpl.ExecuteTemplate(w, "addressEdit.gohtml", ViewData); err != nil {
			log.Fatal.Fatalln(err)
		}
		return
	}
	err = json.Unmarshal(body, &ViewData.Transactions)
	if err != nil {
		log.Error.Println(err)
	}

	if err := tpl.ExecuteTemplate(w, "transaction.gohtml", ViewData); err != nil {
		log.Fatal.Fatalln(err)
	}
}
