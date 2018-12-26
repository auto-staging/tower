package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/types"
)

// GetAllEnvironmentsForRepository gets all Environments where the repository matches name (parameter) from DynamoDB, the received Environments are unmarshaled into
// the array of Environments given in the parameters (call by reference).
// If an error occurs the error gets logged and then returned.
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

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, environments)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetAllEnvironmentsForRepository", "operation": "dynamodb/unmarshalListOfMaps"}, 0)
		return err
	}

	return nil
}

// GetSingleEnvironmentForRepository gets the Environment where repository equals name and branch equals branch from DynamoDB,
// the received Environment gets unmarshaled into the Environment struct given in the parameters (call by reference).
// If an error occurs the error gets logged and then returned.
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

	err = dynamodbattribute.UnmarshalMap(result.Item, environment)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetSingleEnvironmentForRepository", "operation": "dynamodb/unmarshalMap"}, 0)
		return err
	}

	return nil
}

// AddEnvironmentForRepository adds a new Environment for the repository given in the parameters, the values for the new Environment are
// in the EnvironmentPost struct.
// If some values are unset, they will be set with the defaults from the repository.
// After successfully adding the new Environment to DynamoDB, the Builder Lambda gets invoked to add the Schedules and the CodeBuild Job.
// If an error occurs the error gets logged and then returned. If no error occurs the newly created Environment gets returned.
func AddEnvironmentForRepository(environment types.EnvironmentPost, name string) (types.Environment, error) {
	svc := getDynamoDbClient()

	creation := time.Now().UTC()

	inputEnvironment := types.Environment{
		Repository:            name,
		Branch:                environment.Branch,
		Status:                "pending",
		CreationDate:          creation.String(),
		InfrastructureRepoURL: environment.InfrastructureRepoURL,
		ShutdownSchedules:     environment.ShutdownSchedules,
		StartupSchedules:      environment.StartupSchedules,
		EnvironmentVariables:  environment.EnvironmentVariables,
		CodeBuildRoleARN:      environment.CodeBuildRoleARN,
	}

	// Overwrite unset values with defaults from the parent repository
	if inputEnvironment.ShutdownSchedules == nil || inputEnvironment.StartupSchedules == nil || inputEnvironment.EnvironmentVariables == nil || inputEnvironment.InfrastructureRepoURL == "" || inputEnvironment.CodeBuildRoleARN == "" {
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
		if inputEnvironment.InfrastructureRepoURL == "" {
			config.Logger.Log(errors.New("Overwriting InfrastructureRepoURL - Default = "+fmt.Sprint(repository.InfrastructureRepoURL)), map[string]string{"module": "controller/AddEnvironmentForRepository", "operation": "overwrite/InfrastructureRepoURL"}, 4)
			inputEnvironment.InfrastructureRepoURL = repository.InfrastructureRepoURL
		}
		if inputEnvironment.CodeBuildRoleARN == "" {
			config.Logger.Log(errors.New("Overwriting codeBuildRoleARN - Default = "+fmt.Sprint(repository.CodeBuildRoleARN)), map[string]string{"module": "controller/AddEnvironmentForRepository", "operation": "overwrite/codeBuildRoleARN"}, 4)
			inputEnvironment.CodeBuildRoleARN = repository.CodeBuildRoleARN
		}
	}

	av, err := dynamodbattribute.MarshalMap(inputEnvironment)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/AddEnvironmentForRepositroy", "operation": "dynamodb/marshalMap"}, 0)
		return types.Environment{}, err
	}

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

	// Invoke Builder Lambda to configure schedules
	event := types.BuilderEvent{
		Operation:         "UPDATE_SCHEDULE",
		Branch:            inputEnvironment.Branch,
		Repository:        inputEnvironment.Repository,
		ShutdownSchedules: inputEnvironment.ShutdownSchedules,
		StartupSchedules:  inputEnvironment.StartupSchedules,
	}
	body, _ := json.Marshal(event)

	client := getLambdaClient()
	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String("auto-staging-builder"),
		InvocationType: aws.String("Event"),
		Payload:        body,
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/AddEnvironmentForRepositroy", "operation": "builder/invokeSchedule"}, 0)
		return types.Environment{}, err
	}

	// Invoke Builder Lambda to generate environment
	event = types.BuilderEvent{
		Operation:             "CREATE",
		Branch:                inputEnvironment.Branch,
		Repository:            inputEnvironment.Repository,
		CodeBuildRoleARN:      inputEnvironment.CodeBuildRoleARN,
		EnvironmentVariables:  inputEnvironment.EnvironmentVariables,
		InfrastructureRepoURL: inputEnvironment.InfrastructureRepoURL,
	}
	body, _ = json.Marshal(event)

	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String("auto-staging-builder"),
		InvocationType: aws.String("Event"),
		Payload:        body,
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/AddEnvironmentForRepositroy", "operation": "builder/invoke"}, 0)
		return types.Environment{}, err
	}

	return inputEnvironment, nil
}

// UpdateEnvironment updates an existing Environment in DynamoDB where repository equals name and branch equals branch, the updated values for the
// Environment are in the EnvironmentPut struct.
// After successfully updating the Environment in DynamoDB, the Builder Lambda gets invoked to update the Schedules and the CodeBuild Job.
// If an error occurs the error gets logged and then returned. If no error occurs the updated Environment gets returned.
func UpdateEnvironment(environment *types.EnvironmentPut, name string, branch string) (types.Environment, error) {
	svc := getDynamoDbClient()

	updateStruct := types.EnvironmentUpdate{
		InfrastructureRepoURL: environment.InfrastructureRepoURL,
		ShutdownSchedules:     environment.ShutdownSchedules,
		StartupSchedules:      environment.StartupSchedules,
		CodeBuildRoleARN:      environment.CodeBuildRoleARN,
		EnvironmentVariables:  environment.EnvironmentVariables,
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
		UpdateExpression:          aws.String("SET shutdownSchedules = :shutdownSchedules, startupSchedules = :startupSchedules, environmentVariables = :environmentVariables, infrastructureRepoURL = :infrastructureRepoURL, codeBuildRoleARN = :codeBuildRoleARN"),
		ExpressionAttributeValues: update,
		ConditionExpression:       aws.String("attribute_exists(repository) AND attribute_exists(branch)"),
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateEnvironment", "operation": "dynamodb/exec"}, 0)
		return types.Environment{}, err
	}

	// Invoke Builder Lambda to configure schedules
	event := types.BuilderEvent{
		Operation:         "UPDATE_SCHEDULE",
		Branch:            branch,
		Repository:        name,
		ShutdownSchedules: environment.ShutdownSchedules,
		StartupSchedules:  environment.StartupSchedules,
	}
	body, _ := json.Marshal(event)

	client := getLambdaClient()
	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String("auto-staging-builder"),
		InvocationType: aws.String("Event"),
		Payload:        body,
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateSingleRepository", "operation": "builder/invokeSchedule"}, 0)
		return types.Environment{}, err
	}

	// Invoke Builder Lambda to update environment
	event = types.BuilderEvent{
		Operation:             "UPDATE",
		Branch:                branch,
		Repository:            name,
		InfrastructureRepoURL: environment.InfrastructureRepoURL,
		CodeBuildRoleARN:      environment.CodeBuildRoleARN,
		EnvironmentVariables:  environment.EnvironmentVariables,
	}
	body, _ = json.Marshal(event)

	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String("auto-staging-builder"),
		InvocationType: aws.String("Event"),
		Payload:        body,
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateSingleRepository", "operation": "builder/invoke"}, 0)
		return types.Environment{}, err
	}

	response := types.Environment{}
	dynamodbattribute.UnmarshalMap(result.Attributes, &response)

	return response, nil
}

// DeleteSingleEnvironment invokes the Builder Lambda to delete the schedules and the CodeBuild Job with the infrastructure.
// If an error occurs the error gets logged and then returned.
func DeleteSingleEnvironment(environment *types.Environment, name string, branch string) error {
	// Invoke Builder Lambda to delete schedules
	event := types.BuilderEvent{
		Operation:  "DELETE_SCHEDULE",
		Branch:     branch,
		Repository: name,
	}
	body, _ := json.Marshal(event)

	client := getLambdaClient()
	_, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String("auto-staging-builder"),
		InvocationType: aws.String("Event"),
		Payload:        body,
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/DeleteSingleEnvironment", "operation": "builder/invokeSchedule"}, 0)
		return err
	}

	// Invoke Builder Lambda to delete environment
	event = types.BuilderEvent{
		Operation:  "DELETE",
		Branch:     branch,
		Repository: name,
	}
	body, _ = json.Marshal(event)

	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String("auto-staging-builder"),
		InvocationType: aws.String("Event"),
		Payload:        body,
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/DeleteSingleEnvironment", "operation": "builder/invoke"}, 0)
		return err
	}

	return nil
}

// CheckIfEnvironmentsForRepositoryExist checks if Environments in DynamoDB exist where repository equals name from the parameters. If Environments
// were found then true gets returned, otherwise false.
// If an error occurs the error gets logged and then returned.
func CheckIfEnvironmentsForRepositoryExist(name string) (bool, error) {
	svc := getDynamoDbClient()

	result, err := svc.Query(&dynamodb.QueryInput{
		TableName:              aws.String("auto-staging-environments"),
		KeyConditionExpression: aws.String("repository = :repository"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":repository": {
				S: aws.String(name),
			},
		},
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/CheckIfEnvironmentsForRepositoryExist", "operation": "dynamodb/exec"}, 0)
		return false, err
	}

	environments := []types.Environment{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &environments)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/CheckIfEnvironmentsForRepositoryExist", "operation": "dynamodb/unmarshalListOfMaps"}, 0)
		return false, err
	}

	if len(environments) > 0 {
		return true, nil
	}

	return false, nil
}
