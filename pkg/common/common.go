package common

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	MsgServerError  = "500 - Internal Server Error"
	MsgUserNotFound = "404 - User Not found"
)

// GetEnvVar read all vars declared in .env.
func GetEnvVar(v string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(v)
}

// Add two numbers and return the result.
func Add(val1, val2 int) int {
	return val1 + val2
}

// NewNullString sets empty string to sql null value
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NewNullInt64 sets 0 to sql null value
func NewNullInt64(d int) sql.NullInt64 {
	if d == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(d),
		Valid: true,
	}
}

func FetchData(url string) (body []byte, statusCode int, err error) {
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}

// FormatDate parse and format date to YYYY-MM-DD format.
func FormatDate(x string) string {
	td, err := time.Parse("2006-01-02", x)
	if err == nil {
		return td.Format("02-Jan-2006")
	}
	return ""
}

// FormatDateTime parse and format date to YYYY-MM-DD HH:MM:SS format.
func FormatDateTime(x string) string {
	td, err := time.Parse(time.RFC3339, x)
	if err == nil {
		return td.Format("02-Jan-2006 15:04:05")
	}
	return ""
}
