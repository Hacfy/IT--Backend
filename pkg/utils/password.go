package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"unicode"
)

func StrongPasswordValidator(password string) bool {

	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case regexp.MustCompile(`[!@#$%^&*(),.?":{}<>]`).MatchString(string(ch)):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial

}

func GeneratePassword() (string, error) {

	const char = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789[!{@#$%^&*(),.?\":}<>]"
	const passwordLength = 10

	password := make([]byte, passwordLength)

	for i := 0; i < passwordLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(char))))
		if err != nil {
			return "", err
		}

		password[i] += char[n.Int64()]
	}

	return string(password), nil
}
