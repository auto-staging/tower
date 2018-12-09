package controller

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

func TriggerEnvironemtStatusChangeController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	trigger := types.TriggerSchedulePost{}
	err := json.Unmarshal([]byte(request.Body), &trigger)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}

	status := types.EnvironmentStatus{}
	branch, _ := url.PathUnescape(trigger.Branch)
	err = model.GetSingleEnvironmentStatusInformation(&status, trigger.Repository, branch)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	switch trigger.Action {
	case "start":
		if status.Status != "running" && status.Status != "stopped" {
			config.Logger.Log(errors.New("Can't start environment in status = "+status.Status), map[string]string{"module": "controller/TriggerEnvironemtStatusChangeController", "operation": "statusCheck"}, 0)
			return types.InvalidEnvironmentStatusResponse, nil
		}

		result, err := model.TriggerSchedulerLambdaForEnvironment(trigger.Repository, trigger.Branch, "start")
		if err != nil {
			return types.InternalServerErrorResponse, nil
		}
		return events.APIGatewayProxyResponse{Body: result, StatusCode: 200}, nil

	case "stop":
		if status.Status != "running" && status.Status != "stopped" {
			config.Logger.Log(errors.New("Can't stop environment in status = "+status.Status), map[string]string{"module": "controller/TriggerEnvironemtStatusChangeController", "operation": "statusCheck"}, 0)
			return types.InvalidEnvironmentStatusResponse, nil
		}

		result, err := model.TriggerSchedulerLambdaForEnvironment(trigger.Repository, trigger.Branch, "stop")
		if err != nil {
			return types.InternalServerErrorResponse, nil
		}
		return events.APIGatewayProxyResponse{Body: result, StatusCode: 200}, nil

	}

	return types.InvalidRequestBodyResponse, nil
}
