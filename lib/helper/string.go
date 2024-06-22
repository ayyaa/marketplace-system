package helper

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"

	"github.com/labstack/echo/v4"
)

func GenerateAliasName(fullName string, bankCode string) string {
	// Split the full name into separate words
	words := splitName(fullName)

	// Get the first word as the first name
	firstName := ""
	if len(words) > 0 {
		firstName = words[0]
	}

	// Concatenate the first name and bank code
	result := fmt.Sprintf("%s%s", firstName, bankCode)

	return result
}

// Helper function to split a full name into separate words
func splitName(fullName string) []string {
	words := make([]string, 0)

	// Split the full name by spaces
	for _, word := range strings.Split(strings.ToLower(fullName), " ") {
		if word != "" {
			words = append(words, word)
		}
	}

	return words
}

// Get data user from authenticated request context
func GetDataUserFromCtx(ctx echo.Context) (id interface{}) {
	idRaw := ctx.Get("user_id")

	if _, ok := idRaw.(float64); ok {
		idWT := idRaw.(float64)
		id = int64(idWT)
	} else if _, ok := idRaw.(int64); ok {
		id = idRaw.(int64)
	} else if _, ok := idRaw.(int); ok {
		idWT := idRaw.(int)
		id = int64(idWT)
	}

	return
}

func GenerateReferenceNumberString() string {
	b := make([]byte, 15)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:15]
}
