package model

import (
	"os"
	"strconv"

	"github.com/auto-staging/tower/config"
	"github.com/auto-staging/tower/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// GetConfiguration gets the current LogLevel from the env vars and writes the value to the TowerConfiguration struct from the parameters (call by reference).
// If an error occurs the error gets logged and then returned.
func GetConfiguration(configuration *types.TowerConfiguration) error {
	logLevel, err := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	if err != nil {
		config.Logger.Log(err, map[string]string{"module": "model/GetConfiguration"}, 1)
		return err
	}

	configuration.LogLevel = logLevel

	return nil
}

// UpdateConfiguration updates the LogLevel environment variable with the value stored in the TowerConfiguration struct, after calling the AWS update command
// the LogLevel value returned by the command gets stored in the TowerConfiguration struct (overwrites value used in update) from the parameter (call by reference).
// The two values (before and after update) should match.
// If an error occurs the error gets logged and then returned.
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
