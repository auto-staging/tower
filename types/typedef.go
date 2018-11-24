package types

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type TowerConfiguration struct {
	LogLevel int `json:"logLevel"`
}

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Repository struct {
	Repository            string                `json:"repository"`
	InfrastructureRepoURL string                `json:"infrastructureRepoURL"`
	Webhook               bool                  `json:"webhook"`
	Filters               []string              `json:"filters"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
}

type RepositoryUpdate struct {
	InfrastructureRepoURL string                `json:":infrastructureRepoURL"`
	Webhook               bool                  `json:":webhook"`
	Filters               []string              `json:":filters"`
	ShutdownSchedules     []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:":startupSchedules"`
	CodeBuildRoleARN      string                `json:":codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:":environmentVariables"`
}

type EnvironmentGeneralConfig struct {
	ShutdownSchedules    []TimeSchedule        `json:"shutdownSchedules"`
	StartupSchedules     []TimeSchedule        `json:"startupSchedules"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
}

type EnvironmentGeneralConfigUpdate struct {
	ShutdownSchedules    []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules     []TimeSchedule        `json:":startupSchedules"`
	EnvironmentVariables []EnvironmentVariable `json:":environmentVariables"`
}

type Environment struct {
	Repository            string                `json:"repository"`
	Branch                string                `json:"branch"`
	CreationDate          string                `json:"creationDate"`
	Status                string                `json:"status"`
	InfrastructureRepoURL string                `json:"infrastructureRepoURL"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
}

type EnvironmentUpdate struct {
	InfrastructureRepoURL string                `json:":infrastructureRepoURL"`
	ShutdownSchedules     []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:":startupSchedules"`
	CodeBuildRoleARN      string                `json:":codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:":environmentVariables"`
}

type EnvironmentPut struct {
	InfrastructureRepoURL string                `json:"infrastructureRepoURL"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
}

type EnvironmentPost struct {
	Branch                string                `json:"branch"`
	InfrastructureRepoURL string                `json:"infrastructureRepoURL"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
}

type EnvironmentStatus struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Status     string `json:"status"`
}

type Reflector struct {
	Method     string
	Resource   string
	Path       string
	PathParams map[string]string
	Stage      string
	Body       map[string]*json.RawMessage
	Headers    map[string]string
}

type GitHubWebhook struct {
	Ref        string `json:"ref"`
	RefType    string `json:"ref_type"`
	Repository struct {
		Name string `json:"name"`
	}
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
