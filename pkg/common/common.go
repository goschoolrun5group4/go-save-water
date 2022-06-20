package common

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
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
