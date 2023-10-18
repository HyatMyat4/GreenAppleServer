package helper

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)

	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassowrd(user_password string, db_user_password string) (string, bool) {
	err := bcrypt.CompareHashAndPassword([]byte(db_user_password), []byte(user_password))
	check := true
	msg := ""

	if err != nil {
		check = false
		msg = fmt.Sprintf("Password is incorrect")
	}

	return msg, check
}
