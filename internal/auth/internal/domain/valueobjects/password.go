package valueobjects

import (
	"errors"
	"unicode"
)

var (
	ErrWeakPassword = errors.New("password must be at least 8 chars with upper, lower, digit, special")
)

type Password struct {
	plainText string
}

func NewPassword(pwd string) (Password, error) {
	if !isStrongPassword(pwd) {
		return Password{}, ErrWeakPassword
	}
	return Password{plainText: pwd}, nil
}

func (p Password) String() string {
	return p.plainText
}

// optional: expose this for password hashing
func (p Password) Bytes() []byte {
	return []byte(p.plainText)
}

func isStrongPassword(password string) bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
		hasSymbol = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSymbol = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSymbol
}
