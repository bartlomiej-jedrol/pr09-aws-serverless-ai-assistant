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
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	awsInt "github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/aws"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/slack"
)

var (
	ErrorBadRequest          error = errors.New("bad request")
	ErrorInternalServerError error = errors.New("internal server error")
)

// callLambda calls provided lambda with provided body.
func callLambda(lambdaName string, body string) (*lambdaSvc.InvokeOutput, error) {
	log.Printf("INFO: callLambda - calling lambda: %s from router_lambda", lambdaName)

	cfg, err := awsInt.LoadDefaultConfig()
	if err != nil {
		return nil, ErrorInternalServerError
	}

	input := &lambdaSvc.InvokeInput{
		FunctionName: aws.String(lambdaName),
		Payload:      []byte(body),
	}

	lambda := lambdaSvc.NewFromConfig(cfg)
	r, err := lambda.Invoke(context.TODO(), input)
	if err != nil {
		log.Printf("ERROR: callLambda - failed to call lambda: %s, error: %v", lambdaName, err)
		return nil, ErrorInternalServerError
	}
	return r, nil
}

// buildResponseBody builds API Gateway response body.
func buildResponseBody(body any) string {
	log.Printf("INFO: buildResponseBody - building API Gateway response body")

	// In case of error the caller expects response body in the format: "error": "error message"
	if err, ok := body.(error); ok {
		return fmt.Sprintf(`"error": "%v"`, err)
	}
	return ""
}

// BuildAPIResponse builds API Gateway response.
func buildAPIResponse(statusCode int, body any) (*events.APIGatewayProxyResponse, error) {
	log.Printf("INFO: buildAPIResponse - building API Gateway response")

	r := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	r.Body = buildResponseBody(body)
	return r, nil
}

// HandleRequest routes request to handler based on method and availability of "email"
// query parameter.
func HandleRequest(
	request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Logging
	log.Printf("INFO: HandleRequest - handling router_lambda event: %v", request)
	log.Printf("HTTPMethod: %v", request.HTTPMethod)
	log.Printf("Headers: %v", request.Headers)
	log.Printf("PathParameters: %v", request.PathParameters)
	log.Printf("QueryStringParameters: %v", request.QueryStringParameters)
	log.Printf("Body: %v", request.Body)
	log.Println("New logger added for test")

	err := slack.UnmarshalJSON([]byte(request.Body))
	if err != nil {
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	}

	r, err := callLambda("pr09-link-shortener-lambda", request.Body)
	if err != nil {
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	}

	// lambdaOutput := types.LinkShortenerOutputPayload{}
	// err := json.Unmarshal(res.Payload, &lambdaOutput)
	// if err != nil {
	// 	log.Printf("failed to unmarshal JSON: %s", err)
	// }

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(r.Payload),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
