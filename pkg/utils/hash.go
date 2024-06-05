package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(p string) string {
	bData, _ := bcrypt.GenerateFromPassword([]byte(p), 6)

	return string(bData)
}

func CheckPasswordHash(p string, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))

	return err == nil
}
