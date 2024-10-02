package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GenerateToken generates a random token of the specified length
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	token := base64.URLEncoding.EncodeToString(bytes)[:length]
	return token, nil
}

func GetEnvVariable(envVarName string) *string {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file")
	}

	envVar := os.Getenv(envVarName)
	if envVar == "" {
		log.Printf("%s environment variable is blank", envVarName)
		return nil
	}
	return &envVar
}
