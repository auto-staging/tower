package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/types"
)

func GetGlobalRepositoryConfiguration(configuration *types.EnvironmentGeneralConfig, stage string) error {
	svc := getDynamoDbClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("auto-staging-repositories-global-config"),
		Key: map[string]*dynamodb.AttributeValue{
			"stage": {
				S: aws.String(stage),
			},
		},
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetGlobalRepositoryConfiguration", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalMap(result.Item, configuration)

	return nil
}

func UpdateGlobalRepositoryConfiguration(configuration *types.EnvironmentGeneralConfig, stage string) error {
	svc := getDynamoDbClient()

	updateStruct := types.EnvironmentGeneralConfigUpdate{
		ShutdownSchedules:    configuration.ShutdownSchedules,
		StartupSchedules:     configuration.StartupSchedules,
		CodeBuildRoleARN:     configuration.CodeBuildRoleARN,
		EnvironmentVariables: configuration.EnvironmentVariables,
	}

	update, err := dynamodbattribute.MarshalMap(updateStruct)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateGlobalRepositoryConfiguration", "operation": "dynamodb/marshalUpdateMap"}, 0)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("auto-staging-repositories-global-config"),
		Key: map[string]*dynamodb.AttributeValue{
			"stage": {
				S: aws.String(stage),
			},
		},
		UpdateExpression:          aws.String("SET shutdownSchedules = :shutdownSchedules, startupSchedules = :startupSchedules, environmentVariables = :environmentVariables, codeBuildRoleARN = :codeBuildRoleARN"),
		ExpressionAttributeValues: update,
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateGlobalRepositoryConfiguration", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalMap(result.Attributes, configuration)

	return nil
}
