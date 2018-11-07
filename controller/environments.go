package controller

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/janritter/auto-staging-tower/model"
	"gitlab.com/janritter/auto-staging-tower/types"
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

// func PutSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	repository := types.Repository{}
// 	err := json.Unmarshal([]byte(request.Body), &repository)
// 	if err != nil {
// 		config.Logger.Log(err, map[string]string{"module": "controller/PutSingleRepositoryController", "operation": "unmarshal"}, 4)
// 		return types.InvalidRequestBodyResponse, nil
// 	}

// 	err = model.UpdateSingleRepository(&repository, request.PathParameters["name"])

// 	if err != nil {
// 		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
// 			return types.NotFoundErrorResponse, nil
// 		}
// 		return types.InternalServerErrorResponse, nil
// 	}

// 	body, _ := json.Marshal(repository)
// 	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
// }

// func DeleteSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	obj := types.Repository{}
// 	err := model.DeleteSingleRepository(&obj, request.PathParameters["name"])
// 	if err != nil {
// 		return types.InternalServerErrorResponse, nil
// 	}

// 	if obj.Repository == "" {
// 		return types.NotFoundErrorResponse, nil
// 	}

// 	return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
// }
