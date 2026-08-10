package main

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/lambda"
)

func main() {
	handler := lambda.NewLambdaHandler()
	awslambda.Start(handler.HandleRequest)
}
