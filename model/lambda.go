package model

import (
	"os"

	"github.com/auto-staging/tower/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func getLambdaClient() *lambda.Lambda {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/getLambdaClient", "operation": "aws/session"}, 0)
	}

	return lambda.New(sess)
}
