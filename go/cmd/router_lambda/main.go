// Main implements an entry point of the Lambda function.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	awsInt "github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/aws"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/openai"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/slack"
)

var (
	openaiAPIKey             string
	ErrorBadRequest          error = errors.New("bad request")
	ErrorInternalServerError error = errors.New("internal server error")
)

// init initialises execution environment for AWS Lambda.
func init() {
	log.Println("INFO: init - initializing router_lambda")
	apiKey, err := awsInt.GetEnvironmentVariable("OPENAI_API_KEY")
	if err != nil {
		return
	}
	openaiAPIKey = *apiKey
}

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
	resp, err := lambda.Invoke(context.TODO(), input)
	if err != nil {
		log.Printf("ERROR: callLambda - failed to call lambda: %s, error: %v", lambdaName, err)
		return nil, ErrorInternalServerError
	}
	return resp, nil
}

// buildResponseBody builds API Gateway response body.
func buildResponseBody(body any) string {
	switch v := body.(type) {
	case error:
		return fmt.Sprintf(`{"error":"%v"}`, v)
	case string:
		return fmt.Sprintf(`{"response":"%s"}`, v)
	default:
		return ""
	}
}

// BuildAPIResponse builds API Gateway response.
func buildAPIResponse(statusCode int, body any) (*events.APIGatewayProxyResponse, error) {
	log.Printf("INFO: buildAPIResponse - building API Gateway response")

	resp := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp.Body = buildResponseBody(body)
	return resp, nil
}

// extractSkill..
func extractSkill(completion string) (string, error) {
	m := map[string]string{}
	err := json.Unmarshal([]byte(completion), &m)
	if err != nil {
		log.Printf("ERROR: extractSkill - failed to unmarshal completion, %v", err)
		return "", err
	}
	return m["skill"], nil
}

// mapSkillToLambda..
func mapSkillToLambda(skill string) string {
	switch skill {
	case "link_shortener":
		return "pr09-link-shortener-lambda"
	default:
		return ""
	}
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

	slackPayload, err := slack.UnmarshalSlackJSON([]byte(request.Body))
	if err != nil {
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	}

	elements, err := slack.ExtractElements(slackPayload)
	if err != nil {
		return buildAPIResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}

	intention := elements["text"]
	systemContent := `Recognise a skill based on intention provided by user.
		Respond with JSON object. skill (link_shortener|other). 
		Do not add any wrappers like quotes or backticks apart from pure JSON object.
		Examples###
		user: short
		ai: {"skill":"link_shortener"} 
		`
	completion, err := openai.CreateChatCompletion(openaiAPIKey, "gpt-4o", intention, systemContent)
	if err != nil {
		return buildAPIResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}
	log.Printf("completion: %s", completion)

	skill, err := extractSkill(completion)
	log.Printf("skill: %s", skill)
	if err != nil {
		return buildAPIResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}

	lambda := mapSkillToLambda(skill)
	if lambda == "" {
		log.Printf("INFO: HandleRequest - mapping of skill to lambda failed")
		return buildAPIResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}

	// resp, err := callLambda(lambda, request.Body)
	// if err != nil {
	// 	return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	// }

	return buildAPIResponse(http.StatusOK, skill)
	// return buildAPIResponse(http.StatusOK, resp.Payload)
}

func main() {
	lambda.Start(HandleRequest)
}
