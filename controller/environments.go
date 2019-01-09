package controller

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/auto-staging/tower/config"
	"github.com/auto-staging/tower/model"
	"github.com/auto-staging/tower/types"
	"github.com/aws/aws-lambda-go/events"
)

// GetAllEnvironmentsForRepositoryController is the controller function for the GET /repositories/{name}/environments endpoint.
// The "name" path parameter containing the Repository name gets read from the APIGatewayProxyRequest struct
func GetAllEnvironmentsForRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var obj []types.Environment
	err := model.GetAllEnvironmentsForRepository(&obj, request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, err := json.Marshal(obj)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetAllEnvironmentsForRepositoryController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// AddEnvironmentForRepositoryController is the controller function for the POST /repositories/{name}/environments endpoint.
// The "name" path parameter containing the Repository name and the request body containing the information for the new Environment gets read from the APIGatewayProxyRequest struct
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

	body, err := json.Marshal(result)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/AddEnvironmentForRepositoryController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

// GetSingleEnvironmentForRepositoryController is the controller function for the GET /repositories/{name}/environments/{branch} endpoint.
// The "name" path parameter containing the Repository name and the "branch" path parameter containing the branch name gets read from the APIGatewayProxyRequest struct
func GetSingleEnvironmentForRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.Environment{}
	branch, err := url.PathUnescape(request.PathParameters["branch"])
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetSingleEnvironmentForRepositoryController", "operation": "pathUnescape"}, 0)
		return types.InternalServerErrorResponse, nil
	}
	err = model.GetSingleEnvironmentForRepository(&obj, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	body, err := json.Marshal(obj)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetSingleEnvironmentForRepositoryController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// PutSinglEnvironmentForRepositoryController is the controller function for the PUT /repositories/{name}/environments/{branch} endpoint.
// The "name" path parameter containing the Repository name, the "branch" path parameter containing the branch name
// and the request body containing the updated information for the Environment gets read from the APIGatewayProxyRequest struct
func PutSinglEnvironmentForRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status := types.EnvironmentStatus{}
	branch, err := url.PathUnescape(request.PathParameters["branch"])
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSinglEnvironmentForRepositoryController", "operation": "pathUnescape"}, 0)
		return types.InternalServerErrorResponse, nil
	}
	err = model.GetSingleEnvironmentStatusInformation(&status, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if status.Status != "running" && status.Status != "updating failed" {
		config.Logger.Log(errors.New("Can't update environment in status = "+status.Status), map[string]string{"module": "controller/PutSinglEnvironmentForRepositoryController", "operation": "statusCheck"}, 0)
		return types.InvalidEnvironmentStatusResponse, nil
	}

	environment := types.EnvironmentPut{}
	err = json.Unmarshal([]byte(request.Body), &environment)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSinglEnvironmentForRepositoryController", "operation": "unmarshal"}, 4)
		return types.InvalidRequestBodyResponse, nil
	}
	if !validateIAMRoleARN(environment.CodeBuildRoleARN) {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSinglEnvironmentForRepositoryController", "operation": "validateCodeBuildRoleARN"}, 1)
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"codeBuildRoleARN is not a valid IAM Role ARN\" }", StatusCode: 400}, nil
	}

	result, err := model.UpdateEnvironment(&environment, request.PathParameters["name"], branch)

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return types.NotFoundErrorResponse, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, err := json.Marshal(result)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSinglEnvironmentForRepositoryController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// DeleteSingleEnvironmentController is the controller function for the DELETE /repositories/{name}/environments/{branch} endpoint.
// The "name" path parameter containing the Repository name and the "branch" path parameter containing the branch name gets read from the APIGatewayProxyRequest struct
func DeleteSingleEnvironmentController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status := types.EnvironmentStatus{}
	branch, err := url.PathUnescape(request.PathParameters["branch"])
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/DeleteSingleEnvironmentController", "operation": "pathUnescape"}, 0)
		return types.InternalServerErrorResponse, nil
	}
	err = model.GetSingleEnvironmentStatusInformation(&status, request.PathParameters["name"], branch)
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

	err = model.DeleteSingleEnvironment(request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"Invoked Builder\" }", StatusCode: 202}, nil
}
