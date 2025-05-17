package util

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func CapitalizeFirst(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// Helper function to generate a secure random password
func GenerateSecurePassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
