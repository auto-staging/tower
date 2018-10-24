package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetConfigurationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Println(err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("auto-staging-tower-conf"),
		Key: map[string]*dynamodb.AttributeValue{
			"towerStage": {
				S: aws.String(request.RequestContext.Stage),
			},
		},
	})

	if err != nil {
		fmt.Printf("failed to make Query API call, %v", err)
	}

	obj := types.TowerConfiguration{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &obj)
	if err != nil {
		fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func PutConfigurationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Println(err)
	}

	config := types.TowerConfiguration{}
	err = json.Unmarshal([]byte(request.Body), &config)
	if err != nil {
		log.Println(err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":logLevel": {
				N: aws.String(strconv.Itoa(config.LogLevel)),
			},
		},
		TableName: aws.String("auto-staging-tower-conf"),
		Key: map[string]*dynamodb.AttributeValue{
			"towerStage": {
				S: aws.String(request.RequestContext.Stage),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET logLevel = :logLevel"),
	}

	_, err = svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
	}

	body, _ := json.Marshal(config)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
