package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/pkg/utils"
	dub "github.com/dubinc/dub-go"
	"github.com/dubinc/dub-go/models/operations"
)

func ShortenLinkWithDub(token string, longLink string) string {
	d := dub.New(dub.WithSecurity(token))

	req := &operations.CreateLinkRequestBody{
		URL: longLink,
	}

	ctx := context.Background()
	res, err := d.Links.Create(ctx, req)
	if err != nil {
		log.Fatalf("failed to create shorten link with Dub: %s", err)
	}
	return res.ShortLink
}

// lambdaSvc.InvokeOutput
func HandleRequest() {
	token := utils.GetEnvVariable("DUB_API_KEY")
	if token == nil {
		log.Println("failed to get dub api token")
	}

	longLink := "https://gist.github.com/bartlomiej-jedrol/6c1010ae182de054641608b04eecacfe"
	shortLink := ShortenLinkWithDub(*token, longLink)
	log.Println(shortLink)
}

func main() {
	lambda.Start(HandleRequest)
}
