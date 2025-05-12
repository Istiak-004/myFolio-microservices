package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if !isValidEmail(email) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: email}, nil
}

func (e Email) String() string {
	return e.value
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
