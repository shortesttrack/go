package jwt

import (
	"time"
	"st-go/errors"
	"github.com/dgrijalva/jwt-go"
)

type KeyFunc func() ([]byte, error)

type Token struct {
	raw     string
	exp     time.Time
	isValid bool

	claims  map[string]interface{}
	keyFunc KeyFunc
}

func NewToken(string string, keyFunc KeyFunc) *Token {
	return &Token{
		raw: string,
	}
}

func (t *Token) Parse() error {
	if t.keyFunc == nil {
		return errors.New("key func not found")
	}
	key, err := t.keyFunc()
	if err != nil {
		return err
	}
	token, err := jwt.Parse(t.raw, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(key)
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp, ok := claims["exp"].(int64)
		if !ok {
			return errors.New("exp required for token")
		}
		t.exp = time.Unix(exp, 0)
		token.Claims = claims
		t.isValid = token.Valid
		return nil
	}
	return errors.New("claims failure")
}
