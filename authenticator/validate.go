package authenticator

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

var (
	handleRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	phoneRegex  = regexp.MustCompile(`^\+?[0-9]{7,15}$`)
)

func ValidateUserHandle(handle string) error {
	if handle == "" {
		return fmt.Errorf("user_handle is required")
	}
	if !handleRegex.MatchString(handle) {
		return fmt.Errorf("user_handle must be 3-30 characters, alphanumeric or underscore only")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func ValidatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required")
	}
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("phone number must be 7-15 digits, optional leading +")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one digit")
	}
	return nil
}

func ValidateProfileName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("profile name is required")
	}
	if len(name) > 100 {
		return fmt.Errorf("profile name must be 100 characters or less")
	}
	return nil
}
