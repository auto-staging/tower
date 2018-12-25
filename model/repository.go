package model

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/types"
)

// GetAllRepositories reads all Repositories from the DynamoDB Table and unmarshals them into the array of Repository structs from the parameters (call by reference).
// If an error occurs the error gets logged and then returned.
func GetAllRepositories(repositories *[]types.Repository) error {
	svc := getDynamoDbClient()

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("auto-staging-repositories"),
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetAllRepositories", "operation": "dynamodb/exec"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, repositories)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetAllRepositories", "operation": "dynamodb/unmarshalListOfMaps"}, 0)
		return err
	}

	return nil
}

// GetSingleRepository reads a single Repository entry from the DynamoDB Table where repository matches the namen given in the parameters.
// The response gets unmarshaled into the Repository struct from the parameters (call by reference).
// If an error occurs the error gets logged and then returned.
func GetSingleRepository(repository *types.Repository, name string) error {
	svc := getDynamoDbClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("auto-staging-repositories"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
		},
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetSingleRepository", "operation": "dynamodb/exec"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, repository)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetSingleRepository", "operation": "dynamodb/unmarshalListOfMaps"}, 0)
		return err
	}

	return nil
}

// AddRepository adds a new Repository to the DynamoDB Table, the values are received from the Repository struct given in the parameters.
// Some values are overwritten with the global defaults if they were not set, therefore the API Stage is also required since AddRepository calls GetGlobalRepositoryConfiguration
// internaly. To check the stored values, all values in the Repository struct are overwritten with the response of the AWS SDK command (call by reference).
// If an error occurs the error gets logged and then returned.
func AddRepository(repository *types.Repository, stage string) error {
	svc := getDynamoDbClient()

	// Overwrite unset values with general config defaults
	if repository.ShutdownSchedules == nil || repository.StartupSchedules == nil || repository.EnvironmentVariables == nil || repository.CodeBuildRoleARN == "" {
		config.Logger.Log(errors.New("Overwriting unset variables with global defaults"), map[string]string{"module": "model/AddRepository", "operation": "overwrite"}, 4)
		configuration := types.GeneralConfig{}
		err := GetGlobalRepositoryConfiguration(&configuration, stage)
		if err != nil {
			return err
		}
		if repository.ShutdownSchedules == nil {
			config.Logger.Log(errors.New("Overwriting ShutdownSchedules - Default = "+fmt.Sprint(configuration.ShutdownSchedules)), map[string]string{"module": "model/AddRepository", "operation": "overwrite/ShutdownSchedules"}, 4)
			repository.ShutdownSchedules = configuration.ShutdownSchedules
		}
		if repository.StartupSchedules == nil {
			config.Logger.Log(errors.New("Overwriting StartupSchedules - Default = "+fmt.Sprint(configuration.StartupSchedules)), map[string]string{"module": "model/AddRepository", "operation": "overwrite/StartupSchedules"}, 4)
			repository.StartupSchedules = configuration.StartupSchedules
		}
		if repository.EnvironmentVariables == nil {
			config.Logger.Log(errors.New("Overwriting EnvironmentVariables - Default = "+fmt.Sprint(configuration.EnvironmentVariables)), map[string]string{"module": "model/AddRepository", "operation": "overwrite/EnvironmentVariables"}, 4)
			repository.EnvironmentVariables = configuration.EnvironmentVariables
		}
	}

	av, err := dynamodbattribute.MarshalMap(repository)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/AddRepository", "operation": "dynamodb/marshalMap"}, 0)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String("auto-staging-repositories"),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(repository)"),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/AddRepository", "operation": "dynamodb/exec"}, 0)
		return err
	}

	return nil
}

// UpdateSingleRepository updates an existing Repository in DynamoDB where repository matches the given name with the values from the Repository struct
// in the parameters. To check the updated values, all values in the Repository struct are overwritten with the response of the AWS SDK command (call by reference).
// If an error occurs the error gets logged and then returned.
func UpdateSingleRepository(repository *types.Repository, name string) error {
	svc := getDynamoDbClient()

	updateStruct := types.RepositoryUpdate{
		Webhook:               repository.Webhook,
		Filters:               repository.Filters,
		ShutdownSchedules:     repository.ShutdownSchedules,
		StartupSchedules:      repository.StartupSchedules,
		EnvironmentVariables:  repository.EnvironmentVariables,
		InfrastructureRepoURL: repository.InfrastructureRepoURL,
		CodeBuildRoleARN:      repository.CodeBuildRoleARN,
	}

	update, err := dynamodbattribute.MarshalMap(updateStruct)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateSingleRepository", "operation": "dynamodb/marshalUpdateMap"}, 0)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("auto-staging-repositories"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
		},
		UpdateExpression:          aws.String("SET webhook = :webhook, filters = :filters, shutdownSchedules = :shutdownSchedules, startupSchedules = :startupSchedules, environmentVariables = :environmentVariables, infrastructureRepoURL = :infrastructureRepoURL, codeBuildRoleARN = :codeBuildRoleARN"),
		ExpressionAttributeValues: update,
		ConditionExpression:       aws.String("attribute_exists(repository)"),
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateSingleRepository", "operation": "dynamodb/exec"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, repository)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateSingleRepository", "operation": "dynamodb/unmarshalMap"}, 0)
		return err
	}

	return nil
}

// DeleteSingleRepository deletes an existing Repository in DynamoDB where repository matches the given name from the parameters.
// To check the deleted repository, all values in the Repository struct are overwritten with the response of the AWS SDK command (call by reference).
// If an error occurs the error gets logged and then returned.
func DeleteSingleRepository(repository *types.Repository, name string) error {
	svc := getDynamoDbClient()

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("auto-staging-repositories"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/DeleteSingleRepository", "operation": "dynamodb/exec"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, repository)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/DeleteSingleRepository", "operation": "dynamodb/unmarshalMap"}, 0)
		return err
	}

	return nil
}
