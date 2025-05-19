package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hp), nil
}

func CheckPassword(user_password, stored_password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(stored_password), []byte(user_password)); err != nil {
		return err
	}
	return nil
}
