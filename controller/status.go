package controller

import (
	"encoding/json"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/janritter/auto-staging-tower/model"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetAllEnvironmentsStatusInformationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := []types.EnvironmentStatus{}
	err := model.GetAllEnvironmentsStatusInformation(&obj)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func GetSingleEnvironmentStatusInformationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.EnvironmentStatus{}
	branch, _ := url.PathUnescape(request.PathParameters["branch"])
	err := model.GetSingleEnvironmentStatusInformation(&obj, request.PathParameters["name"], branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
