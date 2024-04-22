package util

import "golang.org/x/crypto/bcrypt"

const (
	salt = "sawitpro"
)

func Generate(password string) (string, error) {
	out, err := bcrypt.GenerateFromPassword([]byte(withSalt(password)), bcrypt.DefaultCost)
	return string(out), err
}

func Check(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(withSalt(password)))
	return err == nil
}

func withSalt(password string) string {
	return password + "|" + salt
}
