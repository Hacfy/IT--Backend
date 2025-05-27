package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/labstack/echo/v4"
)

type WarehouseRepo struct {
	db *sql.DB
}

func NewWarehouseRepo(db *sql.DB) *WarehouseRepo {
	return &WarehouseRepo{
		db: db,
	}
}

func (wr *WarehouseRepo) generateUniquePrefix(name string) (string, error) {

	query := database.NewDBinstance(wr.db)

	if len(strings.TrimSpace(name)) == 0 {
		return "", fmt.Errorf("product name cannot be empty")
	}
	name = strings.ToLower(name)

	letters := ""
	seen := make(map[rune]bool)

	for _, ch := range name {
		if unicode.IsLetter(ch) && !seen[ch] {
			letters += string(ch)
			seen[ch] = true
		}
		if len(letters) >= 3 {
			break
		}
	}

	// filler characters
	filler := "zyxwvutsrqponmlkjihgfedcba"

	// pad to length 3
	for len(letters) < 3 {
		letters += string(filler[0])
	}

	prefix := letters
	maxAttempts := 10000
	attempts := 0

	for {
		if !query.IfPrefixExists(prefix) {
			break
		}

		attempts++
		if attempts > maxAttempts {
			return "", fmt.Errorf("failed to generate unique prefix")
		}

		// increment prefix like base-N using filler characters
		runes := []rune(prefix)
		changed := false

		for i := len(runes) - 1; i >= 0; i-- {
			index := strings.IndexRune(filler, runes[i])
			if index == -1 {
				return "", fmt.Errorf("invalid character in prefix: %c", runes[i])
			}
			if index < len(filler)-1 {
				runes[i] = rune(filler[index+1])
				for j := i + 1; j < len(runes); j++ {
					runes[j] = rune(filler[0])
				}
				changed = true
				break
			}
		}

		if !changed {
			// overflow: add one more character to the beginning
			prefix = string(filler[0]) + string(runes)
		} else {
			prefix = string(runes)
		}
	}

	return prefix, nil

}

func (wr *WarehouseRepo) CreateComponent(e echo.Context) (int, error) {
	return http.StatusOK, nil
}

func (wr *WarehouseRepo) DeleteComponent(e echo.Context) (int, error) {
	return http.StatusOK, nil
}
