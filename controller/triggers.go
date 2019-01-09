package controller

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/auto-staging/tower/config"
	"github.com/auto-staging/tower/model"
	"github.com/auto-staging/tower/types"
	"github.com/aws/aws-lambda-go/events"
)

// TriggerEnvironemtStatusChangeController is the controller function for the POST /triggers/schedule endpoint.
// The request body containing the desired status for the Environment gets read from the APIGatewayProxyRequest struct
func TriggerEnvironemtStatusChangeController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	trigger := types.TriggerSchedulePost{}
	err := json.Unmarshal([]byte(request.Body), &trigger)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}

	status := types.EnvironmentStatus{}
	branch, err := url.PathUnescape(trigger.Branch)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/TriggerEnvironemtStatusChangeController", "operation": "pathUnescape"}, 0)
		return types.InternalServerErrorResponse, nil
	}
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
