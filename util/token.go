package util

import (
	"fmt"
	"time"

	"crypto/rsa"
	"io/ioutil"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

const (
	privKeyPath = "rsa"
	pubKeyPath  = "rsa-pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	var err error
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalln(err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalln(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatalln(err)
	}
}

func GenerateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(signKey)
}

func ParseToken(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("no token provided")
	}
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return "", err
	}
	if t.Valid {
		return t.Claims.(jwt.MapClaims)["id"].(string), nil
	}
	return "", nil
}
