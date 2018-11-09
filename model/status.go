package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/janritter/auto-staging-tower/config"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetAllEnvironmentsStatusInformation(status *[]types.EnvironmentStatus) error {
	svc := getDynamoDbClient()

	// TODO Limit DynmoDB query / scan to required attributes
	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("auto-staging-environments"),
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetAllEnvironmentsStatusInformation", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalListOfMaps(result.Items, status)

	return nil
}

func GetSingleEnvironmentStatusInformation(status *types.EnvironmentStatus, name string, branch string) error {
	svc := getDynamoDbClient()

	// TODO Limit DynmoDB query / scan to required attributes
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("auto-staging-environments"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
			"branch": {
				S: aws.String(branch),
			},
		},
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetSingleEnvironmentStatusInformation", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalMap(result.Item, status)

	return nil
}
