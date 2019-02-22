package model

import (
	"os"

	"github.com/auto-staging/tower/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func getDynamoDbClient() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/getDynamoDbClient", "operation": "aws/session"}, 0)
	}

	return dynamodb.New(sess)
}
