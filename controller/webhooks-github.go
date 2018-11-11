package controller

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/janritter/auto-staging-tower/model"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GitHubWebhookPingController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{Body: "{\"message\": \"Pong\"}", StatusCode: 200}, nil
}

func GitHubWebhookCreateController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	webhook := types.GitHubWebhook{}
	err := json.Unmarshal([]byte(request.Body), &webhook)
	if err != nil || webhook.RefType != "branch" {
		return types.InvalidRequestBodyResponse, nil
	}

	repository := types.Repository{}
	model.GetSingleRepository(&repository, webhook.Repository.Name)

	if repository.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	hit := false
	for _, filter := range repository.Filters {
		match, _ := regexp.MatchString(filter, webhook.Ref)
		if match {
			hit = true
			break
		}
	}

	if !hit {
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"No filter match\" }", StatusCode: 400}, nil
	}

	result, err := model.AddEnvironmentForRepository(types.EnvironmentPost{Branch: webhook.Ref}, repository.Repository)
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"Unique constraint violation\" }", StatusCode: 400}, nil
		}
		return types.InternalServerErrorResponse, nil
	}

	body, _ := json.Marshal(result)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}
