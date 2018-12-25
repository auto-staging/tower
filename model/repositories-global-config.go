package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/types"
)

// GetGlobalRepositoryConfiguration reads the current global repository configuration from the DynamoDB Table and unmarshals it to the
// GeneralConfig struct from the parameters (call by refernce).
// Next to the GeneralConfig struct, the stage parameter which is used as Key in DynamoDB and contains the API stage is required.
// If an error occurs the error gets logged and then returned.
func GetGlobalRepositoryConfiguration(configuration *types.GeneralConfig, stage string) error {
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

	err = dynamodbattribute.UnmarshalMap(result.Item, configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetGlobalRepositoryConfiguration", "operation": "dynamodb/unmarshalMap"}, 0)
		return err
	}

	return nil
}

// UpdateGlobalRepositoryConfiguration updates the global repository configuration in DynamoDB by using the AWS SDK with the values
// from the GeneralConfig struct in the parameters, after the update all values in the GeneralConfig struct are overwritten with the AWS command results.
// Next to the GeneralConfig struct, the stage parameter which is used as Key in DynamoDB and contains the API stage is required.
// If an error occurs the error gets logged and then returned.
func UpdateGlobalRepositoryConfiguration(configuration *types.GeneralConfig, stage string) error {
	svc := getDynamoDbClient()

	updateStruct := types.GeneralConfigUpdate{
		ShutdownSchedules:    configuration.ShutdownSchedules,
		StartupSchedules:     configuration.StartupSchedules,
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
		UpdateExpression:          aws.String("SET shutdownSchedules = :shutdownSchedules, startupSchedules = :startupSchedules, environmentVariables = :environmentVariables"),
		ExpressionAttributeValues: update,
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateGlobalRepositoryConfiguration", "operation": "dynamodb/exec"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetGlobalRepositoryConfiguration", "operation": "dynamodb/unmarshalMap"}, 0)
		return err
	}

	return nil
}
