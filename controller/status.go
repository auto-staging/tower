package controller

import (
	"encoding/json"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

// GetAllEnvironmentsStatusInformationController is the controller function for the GET /repositories/environments/status endpoint.
func GetAllEnvironmentsStatusInformationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var obj []types.EnvironmentStatus
	err := model.GetAllEnvironmentsStatusInformation(&obj)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, err := json.Marshal(obj)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetAllEnvironmentsStatusInformationController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// GetSingleEnvironmentStatusInformationController is the controller function for the GET /repositories/{name}/environments/{branch}/status endpoint.
// The "name" path parameter containing the Repository name and the "branch" path parameter containing the branch name gets read from the APIGatewayProxyRequest struct
func GetSingleEnvironmentStatusInformationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.EnvironmentStatus{}
	branch, err := url.PathUnescape(request.PathParameters["branch"])
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetSingleEnvironmentStatusInformationController", "operation": "pathUnescape"}, 0)
		return types.InternalServerErrorResponse, nil
	}
	err = model.GetSingleEnvironmentStatusInformation(&obj, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	body, err := json.Marshal(obj)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetSingleEnvironmentStatusInformationController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
