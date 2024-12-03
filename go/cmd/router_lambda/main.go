// Main implements an entry point of the Lambda function.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/api"
	awsInt "github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/aws"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/openai"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/telegram"
)

var (
	openaiAPIKey             string
	awsConfig                *aws.Config
	ErrorBadRequest          error = errors.New("bad request")
	ErrorInternalServerError error = errors.New("internal server error")
)

// init initialises execution environment for AWS Lambda.
func init() {
	log.Println("INFO: init - initializing router_lambda")
	cfg, err := awsInt.LoadDefaultConfig()
	if err != nil {
		return
	}
	awsConfig = cfg

	apiKey, err := awsInt.GetEnvironmentVariable("OPENAI_API_KEY")
	if err != nil {
		return
	}
	openaiAPIKey = apiKey
}

// callLambda calls provided lambda with provided body.
func callLambda(lambdaName string, body string) (*lambdaSvc.InvokeOutput, error) {
	log.Printf("INFO: callLambda - calling lambda: %s from router_lambda", lambdaName)

	input := &lambdaSvc.InvokeInput{
		FunctionName: aws.String(lambdaName),
		Payload:      []byte(body),
	}

	lambda := lambdaSvc.NewFromConfig(*awsConfig)
	resp, err := lambda.Invoke(context.TODO(), input)
	if err != nil {
		log.Printf("ERROR: callLambda - failed to call lambda: %s, error: %v", lambdaName, err)
		return nil, ErrorInternalServerError
	}
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
	log.Printf("INFO: HandleRequest - handling router_lambda event: %+v", request)
	log.Printf("HTTPMethod: %+v", request.HTTPMethod)
	log.Printf("Headers: %+v", request.Headers)
	log.Printf("PathParameters: %+v", request.PathParameters)
	log.Printf("QueryStringParameters: %+v", request.QueryStringParameters)
	log.Printf("Body: %+v", request.Body)
	log.Println("New logger added for test")

	sourceSystem := request.Headers["Source-System"]
	if sourceSystem != "telegram" {
		return api.BuildResponse(http.StatusBadRequest, ErrorBadRequest)
	}

	message, err := telegram.UnmarshalJSON([]byte(request.Body))
	if err != nil {
		return api.BuildResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}
	log.Printf("telegram message: %+v", message)

	text := message.Text
	log.Printf("telegram text: %+v", text)

	systemContent := `Recognise a skill based on text provided by user.
		Respond with JSON object. skill (link_shortener|other).
		Do not add any wrappers like quotes or backticks apart from pure JSON object.
		Examples###
		user: short ...
		ai: {"skill":"link_shortener"}
		`
	completion, err := openai.CreateChatCompletion(openaiAPIKey, "gpt-4o", text, systemContent)
	if err != nil {
		return api.BuildResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}
	log.Printf("completion: %s", completion)

	skill, err := extractSkill(completion)
	log.Printf("skill: %s", skill)
	if err != nil {
		return api.BuildResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}

	lambda := mapSkillToLambda(skill)
	if lambda == "" {
		log.Printf("INFO: HandleRequest - mapping of skill to lambda failed")
		return api.BuildResponse(http.StatusInternalServerError, ErrorInternalServerError)
	}

	r, err := callLambda("pr09-link-shortener-lambda", text)
	if err != nil {
		return api.BuildResponse(http.StatusBadRequest, ErrorBadRequest)
	}
	log.Printf("r: %+v\n", r)

	return api.BuildResponse(http.StatusOK, "request received")
}

func main() {
	lambda.Start(HandleRequest)
}
