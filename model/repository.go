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

func GetSingleRepository(repository *types.Repository, name string) error {
	svc := getDynamoDbClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("auto-staging-tower-repositories"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
		},
	})

	if err != nil {
		fmt.Printf("failed to make Query API call, %v", err)
	}

	return dynamodbattribute.UnmarshalMap(result.Item, repository)
}

func AddRepository(repository types.Repository) error {
	svc := getDynamoDbClient()

	av, err := dynamodbattribute.MarshalMap(repository)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("auto-staging-tower-repositories"),
		Item:      av,
	}

	_, err = svc.PutItem(input)

	return err
}

func DeleteSingleRepository(repository *types.Repository, name string) error {
	svc := getDynamoDbClient()

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("auto-staging-tower-repositories"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
	})

	dynamodbattribute.UnmarshalMap(result.Attributes, repository)

	return err
}
