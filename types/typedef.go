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
	Repository            string                `json:"repository,omitempty"`
	InfrastructureRepoURL string                `json:"infrastructureRepoURL,omitempty"`
	Webhook               bool                  `json:"webhook,omitempty"`
	Filters               []string              `json:"filters,omitempty"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules,omitempty"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN,omitempty"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables,omitempty"`
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

type GeneralConfig struct {
	ShutdownSchedules    []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules     []TimeSchedule        `json:"startupSchedules,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

type GeneralConfigUpdate struct {
	ShutdownSchedules    []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules     []TimeSchedule        `json:":startupSchedules"`
	EnvironmentVariables []EnvironmentVariable `json:":environmentVariables"`
}

type Environment struct {
	Repository            string                `json:"repository,omitempty"`
	Branch                string                `json:"branch,omitempty"`
	CreationDate          string                `json:"creationDate,omitempty"`
	Status                string                `json:"status,omitempty"`
	InfrastructureRepoURL string                `json:"infrastructureRepoURL,omitempty"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules,omitempty"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN,omitempty"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

type EnvironmentUpdate struct {
	InfrastructureRepoURL string                `json:":infrastructureRepoURL"`
	ShutdownSchedules     []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:":startupSchedules"`
	CodeBuildRoleARN      string                `json:":codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:":environmentVariables"`
}

type EnvironmentPut struct {
	InfrastructureRepoURL string                `json:"infrastructureRepoURL,omitempty"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules,omitempty"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN,omitempty"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

type EnvironmentPost struct {
	Branch                string                `json:"branch,omitempty"`
	InfrastructureRepoURL string                `json:"infrastructureRepoURL,omitempty"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules,omitempty"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN,omitempty"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

type EnvironmentStatus struct {
	Repository string `json:"repository,omitempty"`
	Branch     string `json:"branch,omitempty"`
	Status     string `json:"status,omitempty"`
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

type TriggerSchedulePost struct {
	Branch     string `json:"branch"`
	Repository string `json:"repository"`
	Action     string `json:"action"`
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

var InvalidEnvironmentStatusResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Can't execute operation in current environment status\"}",
	StatusCode: 400,
}

var NotFoundErrorResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Not found\"}",
	StatusCode: 404,
}
