package controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/janritter/auto-staging-tower/model"
	"gitlab.com/janritter/auto-staging-tower/types"
)

func GetConfigurationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj := types.TowerConfiguration{}
	err := model.GetConfiguration(&obj, request.RequestContext.Stage)
	if err != nil {
		fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	body, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func PutConfigurationController(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config := types.TowerConfiguration{}
	err := json.Unmarshal([]byte(request.Body), &config)
	if err != nil {
		log.Println(err)
	}

	err = model.UpdateConfiguration(config, request.RequestContext.Stage)

	if err != nil {
		fmt.Println(err.Error())
	}

	body, _ := json.Marshal(config)

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
