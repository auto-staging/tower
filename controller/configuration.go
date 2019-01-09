package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

// GetConfigurationController is the controller function for the GET /configuration endpoint.
func GetConfigurationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.TowerConfiguration{}
	err := model.GetConfiguration(&obj)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, err := json.Marshal(obj)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetConfigurationController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// PutConfigurationController is the controller function for the PUT /configuration endpoint.
// The request body with the update information gets read from the APIGatewayProxyRequest struct.
func PutConfigurationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	configuration := types.TowerConfiguration{}
	err := json.Unmarshal([]byte(request.Body), &configuration)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}

	err = model.UpdateConfiguration(&configuration)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, err := json.Marshal(configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutConfigurationController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
