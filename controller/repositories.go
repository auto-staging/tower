package controller

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

// GetAllRepositoriesController is the controller function for the GET /repositories endpoint.
func GetAllRepositoriesController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var obj []types.Repository
	err := model.GetAllRepositories(&obj)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// AddRepositoryController is the controller function for the POST /repositories endpoint.
// The request body with the information for the new Repository gets read from the APIGatewayProxyRequest struct.
func AddRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := types.Repository{}
	err := json.Unmarshal([]byte(request.Body), &repo)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}
	if !validateIAMRoleARN(repo.CodeBuildRoleARN) {
		config.Logger.Log(err, map[string]string{"module": "controller/AddRepositoryController", "operation": "validateCodeBuildRoleARN"}, 1)
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"codeBuildRoleARN is not a valid IAM Role ARN\" }", StatusCode: 400}, nil
	}

	err = model.AddRepository(&repo, request.RequestContext.Stage)

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"Unique constraint violation\" }", StatusCode: 400}, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(repo)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

// GetSingleRepositoryController is the controller function for the GET /repositories/{name} endpoint.
// The "name" path parameter gets read from the APIGatewayProxyRequest struct.
func GetSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.Repository{}
	err := model.GetSingleRepository(&obj, request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// PutSingleRepositoryController is the controller function for the PUT /repositories/{name} endpoint.
// The request body containing the information for the new Repository gets read from the APIGatewayProxyRequest struct
func PutSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repository := types.Repository{}
	err := json.Unmarshal([]byte(request.Body), &repository)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSingleRepositoryController", "operation": "unmarshal"}, 4)
		return types.InvalidRequestBodyResponse, nil
	}
	if !validateIAMRoleARN(repository.CodeBuildRoleARN) {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSingleRepositoryController", "operation": "validateCodeBuildRoleARN"}, 1)
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"codeBuildRoleARN is not a valid IAM Role ARN\" }", StatusCode: 400}, nil
	}

	err = model.UpdateSingleRepository(&repository, request.PathParameters["name"])

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return types.NotFoundErrorResponse, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(repository)
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// DeleteSingleRepositoryController is the controller function for the DELETE /repositories/{name} endpoint.
// The "name" path parameter gets read from the APIGatewayProxyRequest struct.
func DeleteSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	exist, err := model.CheckIfEnvironmentsForRepositoryExist(request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if exist {
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"First remove all environments for the repository\" }", StatusCode: 400}, nil
	}

	obj := types.Repository{}
	err = model.DeleteSingleRepository(&obj, request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
}
