package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/controller"
)

// Handler is the main function called by lambda.Start, it redirects the request to the matching controller by resource and http method.
// Since the Lambda function is called through API Gateway it uses APIGatewayProxyRequest as parameter
// to get information about the request (containing ressource, method and much more) and APIGatewayProxyResponse as return value (including http code and response message)
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
		return controller.GetSingleEnvironmentForRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}/environments/{branch}" && request.HTTPMethod == "PUT" {
		return controller.PutSinglEnvironmentForRepositoryController(request)
	}

	if request.Resource == "/repositories/{name}/environments/{branch}" && request.HTTPMethod == "DELETE" {
		return controller.DeleteSingleEnvironmentController(request)
	}

	if request.Resource == "/repositories/environments/status" && request.HTTPMethod == "GET" {
		return controller.GetAllEnvironmentsStatusInformationController(request)
	}

	if request.Resource == "/repositories/{name}/environments/{branch}/status" && request.HTTPMethod == "GET" {
		return controller.GetSingleEnvironmentStatusInformationController(request)
	}

	if request.Resource == "/repositories/environments" && request.HTTPMethod == "GET" {
		return controller.GetGlobalRepositoryConfigController(request)
	}

	if request.Resource == "/repositories/environments" && request.HTTPMethod == "PUT" {
		return controller.PutGlobalRepositoryConfigController(request)
	}

	if request.Resource == "/webhooks/github" && request.HTTPMethod == "POST" && request.Headers["X-GitHub-Event"] == "ping" {
		return controller.GitHubWebhookPingController(request)
	}

	if request.Resource == "/webhooks/github" && request.HTTPMethod == "POST" && request.Headers["X-GitHub-Event"] == "create" {
		return controller.GitHubWebhookCreateController(request)
	}

	if request.Resource == "/webhooks/github" && request.HTTPMethod == "POST" && request.Headers["X-GitHub-Event"] == "delete" {
		return controller.GitHubWebhookDeleteController(request)
	}

	if request.Resource == "/triggers/schedule" && request.HTTPMethod == "POST" {
		return controller.TriggerEnvironemtStatusChangeController(request)
	}

	return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"No controller for requested resource and method found\" }", StatusCode: 400}, nil
}

func main() {
	config.Init()

	lambda.Start(Handler)
}
