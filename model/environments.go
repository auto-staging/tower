package model

import (
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
		Repository:   name,
		Branch:       environment.Branch,
		State:        "pending",
		CreationDate: creation.String(),
	}

	av, err := dynamodbattribute.MarshalMap(inputEnvironment)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("auto-staging-environments"),
		Item:      av,
		//ConditionExpression: aws.String("attribute_not_exists(repository) and attribute_not_exists(branch)"), //TODO Fix check
	}

	_, err = svc.PutItem(input)

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/AddEnvironmentForRepositroy", "operation": "dynamodb/exec"}, 0)
		return types.Environment{}, err
	}

	return inputEnvironment, nil
}
