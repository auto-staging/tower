package controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetAllRepositoriesController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Println(err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("auto-staging-tower-repositories"),
	})

	if err != nil {
		fmt.Printf("failed to make Query API call, %v", err)
	}

	obj := []types.Repository{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &obj)
	if err != nil {
		fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func AddRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Println(err)
	}

	repo := types.Repository{}
	err = json.Unmarshal([]byte(request.Body), &repo)
	if err != nil {
		log.Println(err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(repo)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("auto-staging-tower-repositories"),
		Item:      av,
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	body, _ := json.Marshal(repo)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
