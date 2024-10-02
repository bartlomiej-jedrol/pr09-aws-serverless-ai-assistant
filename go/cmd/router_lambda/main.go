// Main implements an entry point of the Lambda function.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

// HandleRequest routes request to handler based on method and availability of "email"
// query parameter.
func HandleRequest(
	request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Logging
	log.Printf("Request: %v", request)
	log.Printf("HTTPMethod: %v", request.HTTPMethod)
	log.Printf("Headers: %v", request.Headers)
	log.Printf("PathParameters: %v", request.PathParameters)
	log.Printf("QueryStringParameters: %v", request.QueryStringParameters)
	log.Printf("Body: %v", request.Body)

	// Get authorization header
	// authorizationHeader := request.Headers["Authorization"]

	// // Load the default config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	// Create a new Lambda client
	client := lambdaSvc.NewFromConfig(cfg)

	functionName := "pr09-authorizer-lambda"

	// Prepare the Lambda function input
	input := &lambdaSvc.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload: []byte(`{
			"authorizationHeader": ""
		}`),
	}

	// Invoke the Authorizer Lambda function
	r, err := client.Invoke(context.TODO(), input)
	if err != nil {
		log.Printf("failed to invoke Authorizer Lambda function")
	}
	fmt.Println(r)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, World!",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
