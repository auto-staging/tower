package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

func GetAllRepositoriesController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	obj := []types.Repository{}
	err := model.GetAllRepositories(&obj)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func AddRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := types.Repository{}
	err := json.Unmarshal([]byte(request.Body), &repo)
	if err != nil {
		return types.InvalidRequestBodyResponse, nil
	}

	// Overwrite unset values with defaults
	if repo.ShutdownSchedules == nil || repo.StartupSchedules == nil || repo.EnvironmentVariables == nil || repo.CodeBuildRoleARN == "" {
		config.Logger.Log(errors.New("Overwriting unset variables with global defaults"), map[string]string{"module": "controller/AddRepositoryController", "operation": "overwrite"}, 4)
		configuration := types.EnvironmentGeneralConfig{}
		err = model.GetGlobalRepositoryConfiguration(&configuration, request.RequestContext.Stage)
		if err != nil {
			return types.InternalServerErrorResponse, nil
		}
		if repo.ShutdownSchedules == nil {
			config.Logger.Log(errors.New("Overwriting ShutdownSchedules - Default = "+fmt.Sprint(configuration.ShutdownSchedules)), map[string]string{"module": "controller/AddRepositoryController", "operation": "overwrite/ShutdownSchedules"}, 4)
			repo.ShutdownSchedules = configuration.ShutdownSchedules
		}
		if repo.StartupSchedules == nil {
			config.Logger.Log(errors.New("Overwriting StartupSchedules - Default = "+fmt.Sprint(configuration.StartupSchedules)), map[string]string{"module": "controller/AddRepositoryController", "operation": "overwrite/StartupSchedules"}, 4)
			repo.StartupSchedules = configuration.StartupSchedules
		}
		if repo.EnvironmentVariables == nil {
			config.Logger.Log(errors.New("Overwriting EnvironmentVariables - Default = "+fmt.Sprint(configuration.EnvironmentVariables)), map[string]string{"module": "controller/AddRepositoryController", "operation": "overwrite/EnvironmentVariables"}, 4)
			repo.EnvironmentVariables = configuration.EnvironmentVariables
		}
		if repo.CodeBuildRoleARN == "" {
			config.Logger.Log(errors.New("Overwriting codeBuildRoleARN - Default = "+fmt.Sprint(configuration.CodeBuildRoleARN)), map[string]string{"module": "controller/AddRepositoryController", "operation": "overwrite/CodeBuildRoleARN"}, 4)
			repo.CodeBuildRoleARN = configuration.CodeBuildRoleARN
		}
	}

	err = model.AddRepository(repo)

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"Unique constraint violation\" }", StatusCode: 400}, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(repo)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

func GetSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.Repository{}
	err := model.GetSingleRepository(&obj, request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func PutSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repository := types.Repository{}
	err := json.Unmarshal([]byte(request.Body), &repository)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/PutSingleRepositoryController", "operation": "unmarshal"}, 4)
		return types.InvalidRequestBodyResponse, nil
	}

	err = model.UpdateSingleRepository(&repository, request.PathParameters["name"])

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return types.NotFoundErrorResponse, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(repository)
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func DeleteSingleRepositoryController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	exist, err := model.CheckIfEnvironmentsForRepositoryExist(request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if exist {
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"First remove all environments for the repository\" }", StatusCode: 400}, nil
	}

	obj := types.Repository{}
	err = model.DeleteSingleRepository(&obj, request.PathParameters["name"])
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if obj.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
}
