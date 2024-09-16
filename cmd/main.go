// Main implements an entry point of the Lambda function.
package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/handlers"
)

type Test struct {
	Name string `json:"name"`
}

// HandleRequest routes request to handler based on method and availability of "email"
// query parameter.
func HandleRequest(
	request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Logging.
	log.Printf("Request: %v", request)
	log.Printf("HTTPMethod: %v", request.HTTPMethod)
	log.Printf("Headers: %v", request.Headers)
	log.Printf("PathParameters: %v", request.PathParameters)
	log.Printf("QueryStringParameters: %v", request.QueryStringParameters)
	log.Printf("Body: %v", request.Body)

	// Check "email" query parameter existence.
	_, ok := request.QueryStringParameters["email"]

	switch request.HTTPMethod {
	case "POST":
		return handlers.CreateUser(request)
	case "GET":
		// If "email" query parameter provided call GetUser else GetUsers.
		if ok {
			return handlers.GetUser(request)
		} else {
			return handlers.GetUsers()
		}
	case "PUT":
		return handlers.UpdateUser(request)
	case "DELETE":
		return handlers.DeleteUser(request)
	default:
		return handlers.UnhandledHTTPMethod(request)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
