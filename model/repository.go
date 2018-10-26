package model

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetAllRepositories(repositories *[]types.Repository) error {
	svc := getDynamoDbClient()

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("auto-staging-tower-repositories"),
	})

	if err != nil {
		fmt.Printf("failed to make Query API call, %v", err)
	}

	return dynamodbattribute.UnmarshalListOfMaps(result.Items, repositories)
}
