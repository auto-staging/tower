package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

func TriggerEnvironemtStatusChangeController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	trigger := types.TriggerSchedulePost{}
	err := json.Unmarshal([]byte(request.Body), &trigger)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}

	switch trigger.Action {
	case "start":
		result, err := model.TriggerSchedulerLambdaForEnvironment(trigger.Repository, trigger.Branch, "start")
		if err != nil {
			return types.InternalServerErrorResponse, nil
		}
		return events.APIGatewayProxyResponse{Body: result, StatusCode: 200}, nil

	case "stop":
		result, err := model.TriggerSchedulerLambdaForEnvironment(trigger.Repository, trigger.Branch, "stop")
		if err != nil {
			return types.InternalServerErrorResponse, nil
		}
		return events.APIGatewayProxyResponse{Body: result, StatusCode: 200}, nil

	}

	return types.InvalidRequestBodyResponse, nil
}
