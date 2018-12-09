package controller

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

func GetAllEnvironmentsForRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := []types.Environment{}
	err := model.GetAllEnvironmentsForRepository(&obj, request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func AddEnvironmentForRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	env := types.EnvironmentPost{}
	err := json.Unmarshal([]byte(request.Body), &env)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}

	repository := types.Repository{}
	err = model.GetSingleRepository(&repository, request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}
	if repository.Repository == "" {
		return events.APIGatewayProxyResponse{Body: "{\"message\": \"Repository not found\"}", StatusCode: 404}, nil
	}

	result, err := model.AddEnvironmentForRepository(env, request.PathParameters["name"])

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"Unique constraint violation\" }", StatusCode: 400}, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(result)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

func GetSingleEnvironmentForRepository(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.Environment{}
	branch, _ := url.PathUnescape(request.PathParameters["branch"])
	err := model.GetSingleEnvironmentForRepository(&obj, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func PutSinglEnvironmentForRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status := types.EnvironmentStatus{}
	branch, _ := url.PathUnescape(request.PathParameters["branch"])
	err := model.GetSingleEnvironmentStatusInformation(&status, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if status.Status != "running" && status.Status != "stopped" {
		config.Logger.Log(errors.New("Can't update environment in status = "+status.Status), map[string]string{"module": "controller/PutSinglEnvironmentForRepositoryController", "operation": "statusCheck"}, 0)
		return types.InvalidEnvironmentStatusResponse, nil
	}

	environment := types.EnvironmentPut{}
	err = json.Unmarshal([]byte(request.Body), &environment)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSinglEnvironmentForRepositoryController", "operation": "unmarshal"}, 4)
		return types.InvalidRequestBodyResponse, nil
	}

	result, err := model.UpdateEnvironment(&environment, request.PathParameters["name"], branch)

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return types.NotFoundErrorResponse, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func DeleteSingleEnvironmentController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status := types.EnvironmentStatus{}
	branch, _ := url.PathUnescape(request.PathParameters["branch"])
	err := model.GetSingleEnvironmentStatusInformation(&status, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if status.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	if status.Status != "running" && status.Status != "stopped" && status.Status != "initiating failed" && status.Status != "destroying failed" {
		config.Logger.Log(errors.New("Can't delete environment in status = "+status.Status), map[string]string{"module": "controller/DeleteSingleEnvironmentController", "operation": "statusCheck"}, 0)
		return types.InvalidEnvironmentStatusResponse, nil
	}

	env := types.Environment{}
	err = model.DeleteSingleEnvironment(&env, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"Invoked Builder\" }", StatusCode: 202}, nil
}
