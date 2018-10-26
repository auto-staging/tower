package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	lambda.Start(Handler)
}
