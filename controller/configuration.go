package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
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

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// PutConfigurationController is the controller function for the PUT /configuration endpoint.
// The request body with the update information gets read from the APIGatewayProxyRequest struct.
func PutConfigurationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config := types.TowerConfiguration{}
	err := json.Unmarshal([]byte(request.Body), &config)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}

	err = model.UpdateConfiguration(&config)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(config)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
