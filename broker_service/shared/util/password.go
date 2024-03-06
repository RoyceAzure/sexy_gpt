package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hashed password : %w", err)
	}
	return string(hashedPassword), nil
}

// return error  由外部決策  不會回傳明碼
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
