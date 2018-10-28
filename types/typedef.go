package types

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type TowerConfiguration struct {
	LogLevel int `json:"logLevel"`
}

type Repository struct {
	Repository string   `json:"repository"`
	Webhook    bool     `json:"webhook"`
	Filters    []string `json:"filters"`
}

type Reflector struct {
	Method     string
	Resource   string
	Path       string
	PathParams map[string]string
	Stage      string
	Body       map[string]*json.RawMessage
}

var InternalServerErrorResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Internal server error\"}",
	StatusCode: 500,
}
