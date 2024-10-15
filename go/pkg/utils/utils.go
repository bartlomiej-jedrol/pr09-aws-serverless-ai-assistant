package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/pkg/models"
	"github.com/joho/godotenv"
)

// LoadDefaultConfig load default AWS config.
func LoadDefaultConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("ERROR: LoadDefaultConfig - %v, %v", ErrorFailedToLoadDefaultConfig, err)
		return aws.Config{}, ErrorFailedToLoadDefaultConfig
	}
	return cfg, nil
}

// GenerateToken generates a random token of the specified length
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("ERROR: GenerateToken - %v, %v", ErrorFailedToGenerateToken, err)
		return "", ErrorFailedToGenerateToken
	}
	token := base64.URLEncoding.EncodeToString(bytes)[:length]
	return token, nil
}

// GetEnvironmentVariable returns environment variable's value.
func GetEnvironmentVariable(envVarName string) (*string, error) {
	log.Printf("INFO: GetEnvironmentVariable - Entering GetEnvironmentVariable")

	err := godotenv.Load()
	if err != nil {
		log.Printf("ERROR: GetEnvironmentVariable - %v, %v", ErrorFailedToLoadEnvFile, err)
		return nil, ErrorFailedToLoadEnvFile
	}

	ev := os.Getenv(envVarName)
	if ev == "" {
		log.Printf("ERROR: GetEnvironmentVariable - %v, %s", ErrorBlankEnvVar, envVarName)
		return nil, ErrorBlankEnvVar
	}
	return &ev, nil
}

// UnmarshalJSON unmarshals provided JSON.
func UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		log.Printf("ERROR: UnmarshalJSON - %v", ErrorEmptyJSONForUnmarshaling)
		return ErrorEmptyJSONForUnmarshaling
	}

	slackMessages := []models.SlackMessage{}
	err := json.Unmarshal(data, &slackMessages)
	if err != nil {
		log.Printf("ERROR: UnmarshalJSON - %v, %v", ErrorFailedToUnmarshalJSON, err)
		return err
	}
	return nil
}
