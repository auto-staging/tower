package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/janritter/auto-staging-tower/config"
	"gitlab.com/janritter/auto-staging-tower/controller"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.Resource == "/configuration" && request.HTTPMethod == "GET" {
		return controller.GetConfigurationController(request)
	}

	if request.Resource == "/configuration" && request.HTTPMethod == "PUT" {
		return controller.PutConfigurationController(request)
	}

	if request.Resource == "/repositories" && request.HTTPMethod == "GET" {
		return controller.GetAllRepositoriesController(request)
	}

	if request.Resource == "/repositories" && request.HTTPMethod == "POST" {
		return controller.AddRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}" && request.HTTPMethod == "GET" {
		return controller.GetSingleRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}" && request.HTTPMethod == "PUT" {
		return controller.PutSingleRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}" && request.HTTPMethod == "DELETE" {
		return controller.DeleteSingleRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}/environments" && request.HTTPMethod == "GET" {
		return controller.GetAllEnvironmentsForRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}/environments" && request.HTTPMethod == "POST" {
		return controller.AddEnvironmentForRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}/environments/{branch}" && request.HTTPMethod == "GET" {
		return controller.GetSingleEnvironmentForRepository(request)
	}

	if request.Resource == "/repositories/{name}/environments/{branch}" && request.HTTPMethod == "PUT" {
		return controller.PutSinglEnvironmentForRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}/environments/{branch}" && request.HTTPMethod == "DELETE" {
		return controller.DeleteSingleEnvironmentController(request)
	}

	// Default reflector for debugging
	path, _ := url.PathUnescape(request.Path)

	for k := range request.PathParameters {
		unescaped, _ := url.PathUnescape(request.PathParameters[k])
		request.PathParameters[k] = unescaped
	}

	fmt.Println(request)

	var objmap map[string]*json.RawMessage
	json.Unmarshal([]byte(request.Body), &objmap)

	response := &types.Reflector{
		Method:     request.HTTPMethod,
		Resource:   request.Resource,
		Path:       path,
		PathParams: request.PathParameters,
		Stage:      request.RequestContext.Stage,
		Body:       objmap,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	config.Init()

	lambda.Start(Handler)
}
