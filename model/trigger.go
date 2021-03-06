package model

import (
	"strconv"

	"github.com/auto-staging/tower/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// TriggerSchedulerLambdaForEnvironment invokes the Scheduler Lambda Function with the repository, branch and action given in the parameters, action
// can be start or stop.
// If invoking the Scheduler fails the error gets logged and then returned. Otherwise the response message of the Scheduler
// gets unquoted and returned.
func TriggerSchedulerLambdaForEnvironment(repository, branch, action string) (string, error) {
	body := []byte("{ \"repository\": \"" + repository + "\", \"branch\": \"" + branch + "\", \"action\": \"" + action + "\" }")

	client := getLambdaClient()

	response, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String("auto-staging-scheduler"),
		InvocationType: aws.String("RequestResponse"),
		Payload:        body,
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/TriggerSchedulerLambdaForEnvironment", "operation": "scheduler/invoke"}, 0)
		return "", err
	}

	output, err := strconv.Unquote(string(response.Payload))
	config.Logger.Log(err, map[string]string{"module": "model/TriggerSchedulerLambdaForEnvironment", "operation": "strconv/unquote"}, 0)

	if output == "" {
		return "{ \"message\": \"scheduler failed, check the scheduler logs for more information\" }", nil
	}

	return output, nil
}
