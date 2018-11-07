package types

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type TowerConfiguration struct {
	LogLevel int `json:"logLevel"`
}

type Repository struct {
	Repository           string            `json:"repository"`
	Webhook              bool              `json:"webhook"`
	Filters              []string          `json:"filters"`
	ShutdownSchedules    []TimeSchedule    `json:"shutdownSchedules"`
	StartupSchedules     []TimeSchedule    `json:"startupSchedules"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

type RepositoryUpdate struct {
	Webhook              bool              `json:":webhook"`
	Filters              []string          `json:":filters"`
	ShutdownSchedules    []TimeSchedule    `json:":shutdownSchedules"`
	StartupSchedules     []TimeSchedule    `json:":startupSchedules"`
	EnvironmentVariables map[string]string `json:":environmentVariables"`
}

type EnvironmentGeneralConfig struct {
	ShutdownSchedules    []TimeSchedule    `json:"shutdownSchedules"`
	StartupSchedules     []TimeSchedule    `json:"startupSchedules"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

type EnvironmentGeneralConfigUpdate struct {
	ShutdownSchedules    []TimeSchedule    `json:":shutdownSchedules"`
	StartupSchedules     []TimeSchedule    `json:":startupSchedules"`
	EnvironmentVariables map[string]string `json:":environmentVariables"`
}

type Environment struct {
	Repository           string            `json:"repository"`
	Branch               string            `json:"branch"`
	CreationDate         string            `json:"creationDate"`
	State                string            `json:"state"`
	ShutdownSchedules    []TimeSchedule    `json:"shutdownSchedules"`
	StartupSchedules     []TimeSchedule    `json:"startupSchedules"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

type EnvironmentUpdate struct {
	ShutdownSchedules    []TimeSchedule    `json:":shutdownSchedules"`
	StartupSchedules     []TimeSchedule    `json:":startupSchedules"`
	EnvironmentVariables map[string]string `json:":environmentVariables"`
}

type EnvironmentPost struct {
	Branch               string            `json:"branch"`
	ShutdownSchedules    []TimeSchedule    `json:"shutdownSchedules"`
	StartupSchedules     []TimeSchedule    `json:"startupSchedules"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

type Reflector struct {
	Method     string
	Resource   string
	Path       string
	PathParams map[string]string
	Stage      string
	Body       map[string]*json.RawMessage
}

type TimeSchedule struct {
	Cron string `json:"cron"`
}

var InternalServerErrorResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Internal server error\"}",
	StatusCode: 500,
}

var InvalidRequestBodyResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Invalid request body\"}",
	StatusCode: 400,
}

var NotFoundErrorResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Not found\"}",
	StatusCode: 404,
}
