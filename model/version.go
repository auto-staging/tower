package model

import (
	"encoding/json"
	"strconv"

	"github.com/auto-staging/tower/config"
	"github.com/auto-staging/tower/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func GetBuilderVersion(componentVersion *types.SingleComponentVersion) error {
	return getVersionInformationFromAutoStagingLambda(componentVersion, "auto-staging-builder")
}

func GetSchedulerVersion(componentVersion *types.SingleComponentVersion) error {
	return getVersionInformationFromAutoStagingLambda(componentVersion, "auto-staging-scheduler")
}

func getVersionInformationFromAutoStagingLambda(componentVersion *types.SingleComponentVersion, lambdaName string) error {
	event := types.BuilderEvent{
		Operation: "VERSION",
	}
	body, err := json.Marshal(event)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/getVersionInformationFromAutoStagingLambda", "operation": "lambda/marshalEvent"}, 0)
		return err
	}

	client := getLambdaClient()
	result, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(lambdaName),
		InvocationType: aws.String("RequestResponse"),
		Payload:        body,
	})
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/getVersionInformationFromAutoStagingLambda", "operation": "lambda/invoke"}, 0)
		return err
	}

	unquoted, err := strconv.Unquote(string(result.Payload))
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/getVersionInformationFromAutoStagingLambda", "operation": "unquote"}, 0)
		return err
	}

	err = json.Unmarshal([]byte(unquoted), componentVersion)
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/getVersionInformationFromAutoStagingLambda", "operation": "unmarshal"}, 0)
		return err
	}

	return nil
}
