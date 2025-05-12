package valueobjects

import (
	"errors"
	"strings"
)

var (
	ErrEmptyToken = errors.New("token must not be empty")
)

type Token struct {
	value string
}

func NewToken(token string) (Token, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return Token{}, ErrEmptyToken
	}
	return Token{value: token}, nil
}

func (t Token) String() string {
	return t.value
}
