package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/auto-staging/tower/config"
	"github.com/auto-staging/tower/types"
)

// GetAllEnvironmentsStatusInformation reads the repository, branch and status columns from all rows of the environments DynamoDB Table and writes
// them to the Array of EnvironmentStatus structs given in the parameters (call by reference).
// If an error occurs the error gets logged and then returned.
func GetAllEnvironmentsStatusInformation(status *[]types.EnvironmentStatus) error {
	svc := getDynamoDbClient()

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("auto-staging-environments"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"), // Workaround reserved keywoard issue
		},
		ProjectionExpression: aws.String("repository, branch, #status"),
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetAllEnvironmentsStatusInformation", "operation": "dynamodb/exec"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, status)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetAllEnvironmentsStatusInformation", "operation": "dynamodb/unmarshalListOfMaps"}, 0)
		return err
	}

	return nil
}

// GetSingleEnvironmentStatusInformation reads the repository, branch and status columns from the row of the environments DynamoDB Table where
// repository and branch match the values given in the parameters. The result is written to the EnvironmentStatus struct from the parameters (call by reference).
// If an error occurs the error gets logged and then returned.
func GetSingleEnvironmentStatusInformation(status *types.EnvironmentStatus, name string, branch string) error {
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
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"), // Workaround reserved keywoard issue
		},
		ProjectionExpression: aws.String("repository, branch, #status"),
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetSingleEnvironmentStatusInformation", "operation": "dynamodb/exec"}, 0)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, status)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetSingleEnvironmentStatusInformation", "operation": "dynamodb/unmarshalListOfMaps"}, 0)
		return err
	}

	return nil
}
