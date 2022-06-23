package server

import (
	"crypto/tls"
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
