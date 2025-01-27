package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

//
func MakeRefreshToken() (string, error) {
	byteSlice := make([]byte, 32)
	_, err := rand.Read(byteSlice)
	if err != nil {
		return "", fmt.Errorf("Failed to create slice of bytes: %v", err)
	}
	return hex.EncodeToString(byteSlice), nil
}