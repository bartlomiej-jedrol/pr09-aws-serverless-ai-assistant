// Main implements an entry point of the Lambda function.
package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func init() {
	// Load the default config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	// Create a Secrets Manager client
	client := secretsmanager.NewFromConfig(cfg)

	secretName := "pr09-ai-assistant-api-key"

	// Get the secret value
	secretValue, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		log.Fatalf("Failed to get secret value: %v", err)
	}

	// Access the secret value
	secret := aws.ToString(secretValue.SecretString)
	log.Printf("Secret: %s", secret)

}

// HandleRequest handles the request and returns a response.
func HandleRequest(
	request events.APIGatewayCustomAuthorizerRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	log.Printf("Token: %s", token)

	// Remove "Bearer" prefix if it exists
	// token = strings.TrimPrefix(token, "Bearer ")
	return &events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
