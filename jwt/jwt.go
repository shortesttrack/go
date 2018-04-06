package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/am-platform/st-go/errors"
)

type Parser interface {
	Parse(string, func(map[string]interface{}) error) error
	SetPublicKey([]byte)
}

func New(pk []byte) Parser {
	return &parser{pk: pk}
}


type parser struct {
	pk []byte
}

func (p *parser) Parse(tokenString string, handler func(map[string]interface{}) error) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(p.pk)
	})
	if err != nil {
		return errors.Unauthorized
	}
	if token.Valid == false {
		return errors.Unauthorized
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return handler(claims)
	}
	return errors.Unauthorized
}

func (p *parser) SetPublicKey(pk []byte) {
	p.pk = pk
}