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
		keyFunc: keyFunc,
		claims: make(map[string]interface{}),
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
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			return errors.New("exp required for token")
		}
		exp := int64(expFloat)
		t.exp = time.Unix(exp, 0)
		t.claims = claims
		t.isValid = token.Valid
		return nil
	}
	return errors.New("claims failure")
}

func (t *Token) IsValid() bool {
	return t.isValid && t.exp.After(time.Now())
}

func (t *Token) Claims() map[string]interface{} {
	return t.claims
}

func (t *Token) Deadline() time.Time {
	return t.exp
}

func (t *Token) SetRaw(token string) {
	t.raw = token
	t.isValid = false
}

func (t *Token) Invalidate() {
	t.raw = ""
	t.exp = time.Now()
	t.claims = make(map[string]interface{})
	t.isValid = false
}