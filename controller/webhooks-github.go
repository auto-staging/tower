package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/model"
	"gitlab.com/auto-staging/tower/types"
)

// GitHubWebhookPingController is the controller function for the POST /webhooks/github endpoint with X-GitHub-Event = ping.
// GitHub sends the ping event after the Webhook was succesfully added to GitHub.
// The GitHub Webhook endpoint is secured through HMAC.
func GitHubWebhookPingController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: "{\"message\": \"Pong\"}", StatusCode: 200}, nil
}

// GitHubWebhookCreateController is the controller function for the POST /webhooks/github endpoint with X-GitHub-Event = create.
// GitHub sends the create event after a new Git branch was created.
// The GitHub Webhook endpoint is secured through HMAC.
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
	err = model.GetSingleRepository(&repository, webhook.Repository.Name)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	if repository.Repository == "" {
		return types.NotFoundErrorResponse, nil
	}

	hit := false
	for _, filter := range repository.Filters {
		match, err := regexp.MatchString(filter, webhook.Ref)
		if err != nil {
			config.Logger.Log(err, map[string]string{"module": "controller/GitHubWebhookCreateController", "operation": "regexpMatchstring"}, 0)
			return types.InternalServerErrorResponse, nil
		}
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

	body, err := json.Marshal(result)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/GitHubWebhookCreateController", "operation": "marshal"}, 0)
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

// GitHubWebhookDeleteController is the controller function for the POST /webhooks/github endpoint with X-GitHub-Event = delete.
// GitHub sends the delete event after a Git branch was deleted.
// The GitHub Webhook endpoint is secured through HMAC.
func GitHubWebhookDeleteController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !verifyHMAC(request.Body, request.Headers["X-Hub-Signature"]) {
		return events.APIGatewayProxyResponse{Body: "{ \"message\" : \"HMAC validation failed\" }", StatusCode: 400}, nil
	}

	webhook := types.GitHubWebhook{}
	err := json.Unmarshal([]byte(request.Body), &webhook)
	if err != nil || webhook.RefType != "branch" {
		return types.InvalidRequestBodyResponse, nil
	}

	status := types.EnvironmentStatus{}
	err = model.GetSingleEnvironmentStatusInformation(&status, webhook.Repository.Name, webhook.Ref)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}
	if status.Status != "running" && status.Status != "stopped" && status.Status != "initiating failed" && status.Status != "destroying failed" {
		config.Logger.Log(errors.New("Can't delete environment in status = "+status.Status), map[string]string{"module": "controller/GitHubWebhookDeleteController", "operation": "statusCheck"}, 0)
		return types.InvalidEnvironmentStatusResponse, nil
	}

	err = model.DeleteSingleEnvironment(webhook.Repository.Name, webhook.Ref)
	if err != nil {
		return types.InternalServerErrorResponse, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 204}, nil
}

func verifyHMAC(body string, githubHash string) bool {
	messageMAC := githubHash[5:] // first 5 chars are sha1=
	messageMACBuf, err := hex.DecodeString(messageMAC)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/verifyHMAC", "operation": "deocdeString"}, 0)
		return false
	}

	mac := hmac.New(sha1.New, []byte(os.Getenv("WEBHOOK_SECRET_TOKEN")))
	_, err = mac.Write([]byte(body))
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "controller/verifyHMAC", "operation": "deocdeString"}, 0)
		return false
	}
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMACBuf, expectedMAC)
}
