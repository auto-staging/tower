package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/janritter/auto-staging-tower/model"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetAllRepositoriesController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	obj := []types.Repository{}
	err := model.GetAllRepositories(&obj)
	if err != nil {
		fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func AddRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := types.Repository{}
	err := json.Unmarshal([]byte(request.Body), &repo)
	if err != nil {
		log.Println(err)
	}

	err = model.AddRepository(repo)

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"unique constraint violation\" }", StatusCode: 400}, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(repo)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

func GetSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.Repository{}
	err := model.GetSingleRepository(&obj, request.PathParameters["name"])
	if err != nil {
		fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	if obj.Repository == "" {
		return events.APIGatewayProxyResponse{Body: "{}", StatusCode: 404}, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func PutSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repository := types.Repository{}
	err := json.Unmarshal([]byte(request.Body), &repository)
	if err != nil {
		log.Println(err)
	}

	err = model.UpdateSingleRepository(&repository, request.PathParameters["name"])

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return events.APIGatewayProxyResponse{Body: "{}", StatusCode: 404}, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(repository)
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func DeleteSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	obj := types.Repository{}
	err := model.DeleteSingleRepository(&obj, request.PathParameters["name"])
	if err != nil {
		fmt.Printf(err.Error())
	}
	if obj.Repository == "" {
		return events.APIGatewayProxyResponse{Body: "{}", StatusCode: 404}, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
}
