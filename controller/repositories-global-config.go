package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

// GetGlobalRepositoryConfigController is the controller function for the GET /repositories/environments endpoint.
func GetGlobalRepositoryConfigController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.GeneralConfig{}
	err := model.GetGlobalRepositoryConfiguration(&obj, request.RequestContext.Stage)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, err := json.Marshal(obj)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetGlobalRepositoryConfigController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// PutGlobalRepositoryConfigController is the controller function for the PUT /repositories/environments endpoint.
// The request body with the updates information gets read from the APIGatewayProxyRequest struct.
func PutGlobalRepositoryConfigController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	configuration := types.GeneralConfig{}
	err := json.Unmarshal([]byte(request.Body), &configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutGlobalRepositoryConfigController", "operation": "unmarshal"}, 4)
		return types.InvalidRequestBodyResponse, nil
	}

	err = model.UpdateGlobalRepositoryConfiguration(&configuration, request.RequestContext.Stage)

	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, err := json.Marshal(configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutGlobalRepositoryConfigController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
