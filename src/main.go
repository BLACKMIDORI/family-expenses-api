package main

import (
	// This is required to initialize controllers
	_ "family-expenses-api/controllers"
	"github.com/a-h/awsapigatewayv2handler"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
)

func main() {
	if isRunningOnLambda() {
		// Run on AWS Lambda
		httpLambdaHandler := awsapigatewayv2handler.NewLambdaHandler(http.DefaultServeMux)
		lambda.Start(httpLambdaHandler)
	} else {
		// Run local server
		log.Println("Running on 0.0.0.0:8080")
		_ = http.ListenAndServe("0.0.0.0:8080", http.DefaultServeMux)
	}
}

func isRunningOnLambda() bool {
	return os.Getenv("AWS_LAMBDA_RUNTIME_API") != ""
}
