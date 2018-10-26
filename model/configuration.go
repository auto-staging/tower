package model

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetConfiguration(configuration *types.TowerConfiguration, stage string) error {
	svc := getDynamoDbClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("auto-staging-tower-conf"),
		Key: map[string]*dynamodb.AttributeValue{
			"towerStage": {
				S: aws.String(stage),
			},
		},
	})

	if err != nil {
		fmt.Printf("failed to make Query API call, %v", err)
	}

	return dynamodbattribute.UnmarshalMap(result.Item, configuration)
}

func UpdateConfiguration(configuration types.TowerConfiguration, stage string) error {
	svc := getDynamoDbClient()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":logLevel": {
				N: aws.String(strconv.Itoa(configuration.LogLevel)),
			},
		},
		TableName: aws.String("auto-staging-tower-conf"),
		Key: map[string]*dynamodb.AttributeValue{
			"towerStage": {
				S: aws.String(stage),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET logLevel = :logLevel"),
	}

	_, err := svc.UpdateItem(input)

	return err
}
