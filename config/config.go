package config

import (
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	lightning "github.com/janritter/go-lightning-log"
	"gitlab.com/janritter/auto-staging-tower/types"
)

var Logger *lightning.Lightning

func Init() {
	logLevel, _ := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	Logger, _ = lightning.Init(logLevel)
}

func UpdateLambdaConfiguration(configuration types.TowerConfiguration) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Println(err)
	}

	svc := lambda.New(sess)

	_, err = svc.UpdateFunctionConfiguration(&lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String("auto-staging-tower"),
		Environment: &lambda.Environment{
			Variables: map[string]*string{"CONFIGURATION_LOG_LEVEL": aws.String(strconv.Itoa(configuration.LogLevel))},
		},
	})

	if err != nil {
		log.Println(err)
	}
}
