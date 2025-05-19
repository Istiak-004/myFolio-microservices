package valueobjects

import (
	"crypto/rand"
	"encoding/base64"
)

type Token struct {
	TokenString string
}

func NewToken() Token {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return Token{
		TokenString: base64.URLEncoding.EncodeToString(b),
	}
}

func NewTokenWithJTI(jti string) Token {
	return Token{
		TokenString: jti,
	}
}

func (t Token) String() string {
	return string(t.TokenString)
}
func (t Token) IsEmpty() bool {
	return t.String() == ""
}
func (t Token) IsValid() bool {
	return !t.IsEmpty()
}
func (t Token) IsEqual(other Token) bool {
	return t.String() == other.String()
}
func (t Token) IsNotEqual(other Token) bool {
	return !t.IsEqual(other)
}
