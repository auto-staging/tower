package controller

import (
	"encoding/json"

	"github.com/auto-staging/tower/config"
	"github.com/auto-staging/tower/model"
	"github.com/auto-staging/tower/types"
	"github.com/aws/aws-lambda-go/events"
)

func GetVersionsController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	versions := types.ComponentVersions{}

	towerVersion := types.SingleComponentVersion{}
	config.GetVersionInformation(&towerVersion)
	versions.Components = append(versions.Components, towerVersion)

	builderVersion := types.SingleComponentVersion{}
	model.GetBuilderVersion(&builderVersion)
	versions.Components = append(versions.Components, builderVersion)

	schedulerVersion := types.SingleComponentVersion{}
	model.GetSchedulerVersion(&schedulerVersion)
	versions.Components = append(versions.Components, schedulerVersion)

	body, err := json.Marshal(versions)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GetVersionsController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
