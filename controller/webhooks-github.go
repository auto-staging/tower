package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

func GitHubWebhookPingController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{Body: "{\"message\": \"Pong\"}", StatusCode: 200}, nil
}

func GitHubWebhookCreateController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !verifyHMAC(request.Body, request.Headers["X-Hub-Signature"]) {
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"HMAC validation failed\" }", StatusCode: 400}, nil
	}

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

func GitHubWebhookDeleteController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !verifyHMAC(request.Body, request.Headers["X-Hub-Signature"]) {
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"HMAC validation failed\" }", StatusCode: 400}, nil
	}

	webhook := types.GitHubWebhook{}
	err := json.Unmarshal([]byte(request.Body), &webhook)
	if err != nil || webhook.RefType != "branch" {
		return types.InvalidRequestBodyResponse, nil
	}

	environment := types.Environment{}
	err = model.DeleteSingleEnvironment(&environment, webhook.Repository.Name, webhook.Ref)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if environment.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
}

func verifyHMAC(body string, githubHash string) bool {
	messageMAC := githubHash[5:] // first 5 chars are sha1=
	messageMACBuf, _ := hex.DecodeString(messageMAC)

	mac := hmac.New(sha1.New, []byte(os.Getenv("WEBHOOK_SECRET_TOKEN")))
	mac.Write([]byte(body))
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMACBuf, expectedMAC)
}
