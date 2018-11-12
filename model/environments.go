package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/janritter/auto-staging-tower/config"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetAllEnvironmentsForRepository(environments *[]types.Environment, name string) error {
	svc := getDynamoDbClient()

	result, err := svc.Query(&dynamodb.QueryInput{
		TableName:              aws.String("auto-staging-environments"),
		KeyConditionExpression: aws.String("#repository = :repository"),
		ExpressionAttributeNames: map[string]*string{
			"#repository": aws.String("repository"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":repository": {
				S: aws.String(name),
			},
		},
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetAllEnvironmentsForRepository", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalListOfMaps(result.Items, environments)

	return nil
}

func GetSingleEnvironmentForRepository(environment *types.Environment, name string, branch string) error {
	svc := getDynamoDbClient()

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
		config.Logger.Log(err, map[string]string{"module": "model/GetSingleEnvironmentForRepository", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalMap(result.Item, environment)

	return nil
}

func AddEnvironmentForRepository(environment types.EnvironmentPost, name string) (types.Environment, error) {
	svc := getDynamoDbClient()

	creation := time.Now().UTC()

	inputEnvironment := types.Environment{
		Repository:           name,
		Branch:               environment.Branch,
		Status:               "pending",
		CreationDate:         creation.String(),
		ShutdownSchedules:    environment.ShutdownSchedules,
		StartupSchedules:     environment.StartupSchedules,
		EnvironmentVariables: environment.EnvironmentVariables,
	}

	// Overwrite unset values with defaults from the parent repository
	if inputEnvironment.ShutdownSchedules == nil || inputEnvironment.StartupSchedules == nil || inputEnvironment.EnvironmentVariables == nil {
		config.Logger.Log(errors.New("Overwriting unset variables with global defaults"), map[string]string{"module": "model/AddEnvironmentForRepository", "operation": "overwrite"}, 4)
		repository := types.Repository{}
		err := GetSingleRepository(&repository, name)
		if err != nil {
			return types.Environment{}, err
		}
		if inputEnvironment.ShutdownSchedules == nil {
			config.Logger.Log(errors.New("Overwriting ShutdownSchedules - Default = "+fmt.Sprint(repository.ShutdownSchedules)), map[string]string{"module": "controller/AddEnvironmentForRepository", "operation": "overwrite/ShutdownSchedules"}, 4)
			inputEnvironment.ShutdownSchedules = repository.ShutdownSchedules
		}
		if inputEnvironment.StartupSchedules == nil {
			config.Logger.Log(errors.New("Overwriting StartupSchedules - Default = "+fmt.Sprint(repository.StartupSchedules)), map[string]string{"module": "controller/AddEnvironmentForRepository", "operation": "overwrite/StartupSchedules"}, 4)
			inputEnvironment.StartupSchedules = repository.StartupSchedules
		}
		if inputEnvironment.EnvironmentVariables == nil {
			config.Logger.Log(errors.New("Overwriting EnvironmentVariables - Default = "+fmt.Sprint(repository.EnvironmentVariables)), map[string]string{"module": "controller/AddEnvironmentForRepository", "operation": "overwrite/EnvironmentVariables"}, 4)
			inputEnvironment.EnvironmentVariables = repository.EnvironmentVariables
		}
	}

	av, err := dynamodbattribute.MarshalMap(inputEnvironment)

	input := &dynamodb.PutItemInput{
		TableName:           aws.String("auto-staging-environments"),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(repository) AND attribute_not_exists(branch)"),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/AddEnvironmentForRepositroy", "operation": "dynamodb/exec"}, 0)
		return types.Environment{}, err
	}

	return inputEnvironment, nil
}

func UpdateEnvironment(environment *types.EnvironmentPut, name string, branch string) (types.Environment, error) {
	svc := getDynamoDbClient()

	updateStruct := types.EnvironmentUpdate{
		ShutdownSchedules:    environment.ShutdownSchedules,
		StartupSchedules:     environment.StartupSchedules,
		EnvironmentVariables: environment.EnvironmentVariables,
	}

	update, err := dynamodbattribute.MarshalMap(updateStruct)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateSingleRepository", "operation": "dynamodb/marshalUpdateMap"}, 0)
		return types.Environment{}, err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("auto-staging-environments"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
			"branch": {
				S: aws.String(branch),
			},
		},
		UpdateExpression:          aws.String("SET shutdownSchedules = :shutdownSchedules, startupSchedules = :startupSchedules, environmentVariables = :environmentVariables"),
		ExpressionAttributeValues: update,
		ConditionExpression:       aws.String("attribute_exists(repository) AND attribute_exists(branch)"),
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateEnvironment", "operation": "dynamodb/exec"}, 0)
		return types.Environment{}, err
	}

	response := types.Environment{}
	dynamodbattribute.UnmarshalMap(result.Attributes, &response)

	return response, nil
}

func DeleteSingleEnvironment(environment *types.Environment, name string, branch string) error {
	svc := getDynamoDbClient()

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("auto-staging-environments"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(name),
			},
			"branch": {
				S: aws.String(branch),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/DeleteSingleEnvironment", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalMap(result.Attributes, environment)

	return nil
}
