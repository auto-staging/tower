package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/auto-staging/tower/config"
)

func getLambdaClient() *lambda.Lambda {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/getLambdaClient", "operation": "aws/session"}, 0)
	}

	return lambda.New(sess)
}
