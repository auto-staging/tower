package model

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/janritter/auto-staging-tower/config"
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
		config.Logger.Log(err, map[string]string{"module": "model/GetConfiguration", "operation": "db/execution"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetConfiguration", "operation": "db/unmarshal"}, 1)
		return err
	}

	return nil
}

func UpdateConfiguration(configuration *types.TowerConfiguration, stage string) error {
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

	result, err := svc.UpdateItem(input)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateConfiguration", "operation": "db/execution"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateConfiguration", "operation": "db/unmarshal"}, 1)
		return err
	}

	return err
}
