package model

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"gitlab.com/auto-staging/tower/config"
	"gitlab.com/auto-staging/tower/types"
)

func GetConfiguration(configuration *types.TowerConfiguration) error {
	logLevel, err := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetConfiguration"}, 1)
		return err
	}

	configuration.LogLevel = logLevel

	return nil
}

func UpdateConfiguration(configuration *types.TowerConfiguration) error {
	sess, errSession := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if errSession != nil {
		config.Logger.Log(errSession, map[string]string{"module": "model/UpdateConfiguration", "operation": "aws/session"}, 0)
		return errSession
	}

	svc := lambda.New(sess)

	result, err := svc.UpdateFunctionConfiguration(&lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String("auto-staging-tower"),
		Environment: &lambda.Environment{
			Variables: map[string]*string{"CONFIGURATION_LOG_LEVEL": aws.String(strconv.Itoa(configuration.LogLevel))},
		},
	})

	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/UpdateConfiguration", "operation": "lambda/update_config"}, 0)
		return err
	}

	logLevel, err := strconv.Atoi(*result.Environment.Variables["CONFIGURATION_LOG_LEVEL"])
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetConfiguration", "operation": "lambda/update_config_result"}, 1)
		return err
	}
	configuration.LogLevel = logLevel

	return nil
}
