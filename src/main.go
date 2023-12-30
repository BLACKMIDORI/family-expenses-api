package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	// This is required to initialize controllers
	_ "family-expenses-api/controllers"
	"github.com/a-h/awsapigatewayv2handler"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
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

func teste() {
	refreshToken := ""
	for i := 0; i < 10; i++ {
		sha_512 := sha512.New()
		sha_512.Write([]byte(uuid.New().String()))
		refreshToken += hex.EncodeToString(sha_512.Sum(nil))
	}
	fmt.Println(len(refreshToken), refreshToken)
}
