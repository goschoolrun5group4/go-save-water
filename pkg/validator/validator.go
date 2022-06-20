// Package validator implements validation functions to validate user inputs.
package validator

import (
	"net/mail"
	"regexp"
	"unicode"
)

// const of different regex and field size.
const (
	name              = "^[a-zA-Z_. ]*$"
	username          = "^[a-zA-Z0-9][a-zA-Z0-9\\_\\-\\.]*[a-zA-Z0-9]$"
	usernameMinLength = 5
	usernameMaxLength = 45
	passwordMinLength = 7
	emailMaxLength    = 100
)

// IsEmpty checks if user input is empty.
func IsEmpty(input string) bool {
	if input == "" && len(input) == 0 {
		return true
	}
	return false
}

// IsValidName validate name against name regex.
// Name consists only of alphabet (a-zA-Z), dot (.), underscore (_) and space( ).
func IsValidName(input string) bool {
	if len(input) > usernameMaxLength {
		return false
	}
	regex := regexp.MustCompile(name)
	return regex.MatchString(input)
}

// IsValidUsername validate username against username regex.
// Username consists between 5 and 20 characters.
// Username consists of alphanumeric characters (a-zA-Z0-9), lowercase, or uppercase.
// Username allowed of the dot (.), underscore (_), and hyphen (-).
// The dot (.), underscore (_), or hyphen (-) cannot be the first or last character.
func IsValidUsername(input string) bool {
	if len(input) < usernameMinLength || len(input) > usernameMaxLength {
		return false
	}
	regex := regexp.MustCompile(username)
	return regex.MatchString(input)
}

// IsValidPassword validate password against unicode package.
// Password has a minimum length of 7 characters.
// Password consist of at least 1 upper and lower case.
// Password consist of at least 1 special character
func IsValidPassword(input string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(input) >= passwordMinLength {
		hasMinLen = true
	}
	for _, char := range input {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// IsValidEmail email address
func IsValidEmail(input string) bool {
	if len(input) > emailMaxLength {
		return false
	}
	_, err := mail.ParseAddress(input)
	return err == nil
}
