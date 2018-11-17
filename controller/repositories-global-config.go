package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

func GetGlobalRepositoryConfigController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.EnvironmentGeneralConfig{}
	err := model.GetGlobalRepositoryConfiguration(&obj, request.RequestContext.Stage)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func PutGlobalRepositoryConfigController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	configuration := types.EnvironmentGeneralConfig{}
	err := json.Unmarshal([]byte(request.Body), &configuration)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutGlobalRepositoryConfigController", "operation": "unmarshal"}, 4)
		return types.InvalidRequestBodyResponse, nil
	}

	err = model.UpdateGlobalRepositoryConfiguration(&configuration, request.RequestContext.Stage)

	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(configuration)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
