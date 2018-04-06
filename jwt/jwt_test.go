package jwt

import (
	"testing"
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/am-platform/st-go/errors"
)

func TestParser_Parse(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := pem.EncodeToMemory(&pem.Block{
		Type:  `RSA PUBLIC KEY`,
		Bytes: bytes,
	})

	parser := New(publicKey)

	testKey := "testkey"
	testValue := "testvalue"

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		testKey: testValue,
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	err = parser.Parse(tokenString, func(claims map[string]interface{}) error {
		if resultValue, ok := claims[testKey].(string); ok {
			if resultValue != testValue {
				return errors.New("token values didnt match")
			}
			return nil
		} else {
			return errors.New(`value not found`)
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestParser_SetPublicKey(t *testing.T) {

}
