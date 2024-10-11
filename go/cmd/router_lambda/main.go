// Main implements an entry point of the Lambda function.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/pkg/utils"
)

var (
	ErrorBadRequest = errors.New("bad request")
)

func callLambda(lambdaName string, body string) (*lambdaSvc.InvokeOutput, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load default config: %v", err)
		return nil, errors.New("failed to load default config")
	}
	log.Printf("CallLambda lambdaName: %s", lambdaName)
	log.Printf("CallLambda body: %s", body)

	input := &lambdaSvc.InvokeInput{
		FunctionName: aws.String(lambdaName),
		Payload:      []byte(body),
	}
	log.Printf("input of router_lambda to target lambda - input: FunctionName: %s", *input.FunctionName)
	log.Printf("input of router_lambda to target lambda - input: Payload: %s", string(input.Payload))

	lambda := lambdaSvc.NewFromConfig(cfg)
	res, err := lambda.Invoke(context.TODO(), input)
	log.Printf("response from target lambda (CallLambda): %v", *res)
	log.Printf("response from target lambda (CallLambda) - Payload: %v", string(res.Payload))
	log.Printf("response from target lambda (CallLambda) - StatusCode: %v", res.StatusCode)
	if err != nil {
		log.Printf("failed to invoke Link Shortener Lambda function: %s", err)
	}
	return res, nil
}

// buildResponseBody builds API Gateway response body
func buildResponseBody(body any) string {
	if err, ok := body.(error); ok {
		return fmt.Sprintf(`"error": "%s"`, err.Error())
	}
	return ""
}

// BuildAPIResponse builds API Gateway response.
func buildAPIResponse(statusCode int, body any) (*events.APIGatewayProxyResponse, error) {
	response := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response.Body = buildResponseBody(body)
	return response, nil
}

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

	err := utils.UnmarshalJSON([]byte(request.Body))
	if err != nil {
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	}

	res, _ := callLambda("pr09-link-shortener-lambda", request.Body)
	log.Printf("response from target lambda (HandleRequest): %v", *res)

	// lambdaOutput := types.LinkShortenerOutputPayload{}
	// err := json.Unmarshal(res.Payload, &lambdaOutput)
	// if err != nil {
	// 	log.Printf("failed to unmarshal JSON: %s", err)
	// }

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(res.Payload),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
