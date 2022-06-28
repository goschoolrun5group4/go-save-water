package server

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	com "go-save-water/pkg/common"
	"go-save-water/pkg/log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	gomail "gopkg.in/mail.v2"
)

type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Data struct {
	Month string `json:"x,omitempty"`
	Usage string `json:"y,omitempty"`
}

type UserUsage struct {
	AccountNumber int    `json:"accountNumber"`
	BillDate      string `json:"billDate"`
	Consumption   string `json:"consumption"`
}

func authenticationCheck(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	var userInfo map[string]interface{}
	cookie, err := r.Cookie(com.GetEnvVar("COOKIE_NAME"))
	if err != nil {
		return nil, err
	}
	url := com.GetEnvVar("API_AUTHENTICATION_ADDR") + "/verify/" + cookie.Value
	body, statusCode, err := com.FetchData(url)
	if err != nil || statusCode == http.StatusNotFound {
		cookie.Expires = time.Now()
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil, err
	}
	json.Unmarshal(body, &userInfo)
	return userInfo, nil
}

// createNewSecureCookie creates and return a new secure cookie.
func createNewSecureCookie(uuid string, expireDT time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     com.GetEnvVar("COOKIE_NAME"),
		Expires:  expireDT,
		Value:    uuid,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
	}
	return cookie
}

func generateJWT(email string) (string, error) {

	jwtKey := []byte(com.GetEnvVar("JWT_SECRET"))

	// Create the Claims
	claims := JwtClaims{
		email,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (bool, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(com.GetEnvVar("JWT_SECRET")), nil
	})

	claim := token.Claims.(*JwtClaims)

	if token.Valid {
		return true, claim.Email, nil
	} else {
		log.Error.Println(err)
		return false, claim.Email, err
	}
}

func sendVerificationEmail(email string) {
	m := gomail.NewMessage()
	m.SetHeader("From", com.GetEnvVar("EMAIL"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Go Save Water - Verify your account")
	// Set E-Mail body. You can set plain text or html with text/html

	token, err := generateJWT(email)
	if err != nil {
		log.Error.Println(err)
	}

	url := "http://localhost:8080/verification/" + token
	body := "<html>" +
		"<div>Thanks for signing up!</div><br/>" +
		"<div>Your account has been created, you can activate your account by pressing the url below. The link will expires in 15 minutes.</div><br/><br/>" +
		"<a href=" + url + ">" + url + "</a>" +
		"</html>"
	m.SetBody("text/html", body)
	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, com.GetEnvVar("EMAIL"), com.GetEnvVar("EMAIL_PASSWORD"))
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getUserUsage(accountNum string, chn chan string) {

	var (
		retJSON []map[string]interface{}
		cData   []Data
	)

	url := com.GetEnvVar("API_USAGE_ADDR") + fmt.Sprintf("/usage/user/%s/latest/6", accountNum)
	body, _, err := com.FetchData(url)
	if err != nil {
		log.Error.Println(err)
		chn <- ""
	}

	err = json.Unmarshal(body, &retJSON)
	if err != nil {
		log.Error.Println(err)
		chn <- ""
	}

	for _, v := range retJSON {
		var pData Data
		t, err := time.Parse("2006-01-02", v["billDate"].(string))
		if err != nil {
			log.Error.Println("Error while parsing date :", err)
		}
		pData.Month = t.Month().String()
		pData.Usage = v["consumption"].(string)
		cData = append(cData, pData)
	}

	jString, err := json.Marshal(cData)
	if err != nil {
		log.Error.Println(err)
		chn <- ""
	} else {
		chn <- string(jString)
	}
}

func getNationalUsage(chn chan string) {
	var (
		retJSON []map[string]interface{}
		cData   []Data
	)

	url := com.GetEnvVar("API_USAGE_ADDR") + "/usage/national/latest/6"
	body, _, err := com.FetchData(url)

	if err != nil {
		log.Error.Println(err)
		chn <- ""
	}

	json.Unmarshal(body, &retJSON)

	for _, v := range retJSON {
		var pData Data
		t, err := time.Parse("2006-01", v["billDate"].(string))
		if err != nil {
			log.Error.Println("Error while parsing date :", err)
		}
		pData.Month = t.Month().String()
		pData.Usage = v["consumption"].(string)
		cData = append(cData, pData)
	}

	jString, err := json.Marshal(cData)
	if err != nil {
		log.Error.Println(err)
		chn <- ""
	} else {
		chn <- string(jString)
	}
}

func postToApi(url string, jsonStr []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}
	return resp, nil
}

func searchIfDateExist(arr []UserUsage, inputDateStr string) bool {
	inputDate, err := time.Parse("2006-01-02", inputDateStr)
	if err != nil {
		return false
	}
	inputYearMonth := inputDate.Format("2006-01")
	return recursiveBinarySearch(len(arr), 0, len(arr)-1, arr, inputYearMonth)
}

// Recursive Binary Search
func recursiveBinarySearch(n, first, last int, arr []UserUsage, inputYearMonth string) bool {
	if first > last {
		return false
	} else {
		mid := (first + last) / 2
		midValue := arr[mid]
		midValueDate, err := time.Parse("2006-01-02", midValue.BillDate)
		if err != nil {
			return false
		}
		midValueYearMonth := midValueDate.Format("2006-01")
		if inputYearMonth == midValueYearMonth {
			return true
		} else {
			if inputYearMonth > midValueYearMonth {
				return recursiveBinarySearch(n, first, mid-1, arr, inputYearMonth)
			} else {
				return recursiveBinarySearch(n, mid+1, last, arr, inputYearMonth)
			}
		}
	}
}

func getReward(rewardID string, chn chan map[string]interface{}) {
	var retJSON map[string]interface{}

	url := com.GetEnvVar("API_REWARD_ADDR") + fmt.Sprintf("/reward/%s", rewardID)
	body, _, err := com.FetchData(url)

	if err != nil {
		log.Error.Println(err)
		chn <- nil
		return
	}

	err = json.Unmarshal(body, &retJSON)

	if err != nil {
		log.Error.Println(err)
		chn <- nil
		return
	}

	chn <- retJSON
}

func getRewards(chn chan []map[string]interface{}) {
	var (
		retJSON []map[string]interface{}
	)

	url := com.GetEnvVar("API_REWARD_ADDR") + "/rewards"
	body, _, err := com.FetchData(url)

	if err != nil {
		log.Error.Println(err)
		chn <- nil
		return
	}

	err = json.Unmarshal(body, &retJSON)

	if err != nil {
		log.Error.Println(err)
		chn <- nil
		return
	}

	chn <- retJSON
}

func removeCurrReward(arr []map[string]interface{}, rewardID string) []map[string]interface{} {
	for k, v := range arr {
		if rewardID == v["rewardID"].(string) {
			return RemoveIndex(arr, k)
		}
	}
	return arr
}

func RemoveIndex(arr []map[string]interface{}, index int) []map[string]interface{} {
	return append(arr[:index], arr[index+1:]...)
}
