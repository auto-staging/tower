package types

import (
	"github.com/aws/aws-lambda-go/events"
)

// TowerConfiguration is the implementation of the TowerAPI TowerConfiguration schema
type TowerConfiguration struct {
	LogLevel           int    `json:"logLevel"`
	WebhookSecretToken string `json:"webhookSecretToken"`
}

// EnvironmentVariable is the implementation of the TowerAPI EnvironmentVariable schema
type EnvironmentVariable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Repository is the implementation of the TowerAPI Repository schema
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

// RepositoryUpdate struct is used for DynamoDB updates, because the update command requires all json keys to start with ":"
type RepositoryUpdate struct {
	InfrastructureRepoURL string                `json:":infrastructureRepoURL"`
	Webhook               bool                  `json:":webhook"`
	Filters               []string              `json:":filters"`
	ShutdownSchedules     []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:":startupSchedules"`
	CodeBuildRoleARN      string                `json:":codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:":environmentVariables"`
}

// GeneralConfig is the implementation of the TowerAPI GeneralConfiguration schema
type GeneralConfig struct {
	ShutdownSchedules    []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules     []TimeSchedule        `json:"startupSchedules,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

// GeneralConfigUpdate struct is used for DynamoDB updates, because the update command requires all json keys to start with ":"
type GeneralConfigUpdate struct {
	ShutdownSchedules    []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules     []TimeSchedule        `json:":startupSchedules"`
	EnvironmentVariables []EnvironmentVariable `json:":environmentVariables"`
}

// Environment is the implementation of the TowerAPI Environment schema
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

// EnvironmentUpdate struct is used for DynamoDB updates, because the update command requires all json keys to start with ":"
type EnvironmentUpdate struct {
	InfrastructureRepoURL string                `json:":infrastructureRepoURL"`
	ShutdownSchedules     []TimeSchedule        `json:":shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:":startupSchedules"`
	CodeBuildRoleARN      string                `json:":codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:":environmentVariables"`
}

// EnvironmentPut is the implementation of the TowerAPI EnvironmentPutBody schema
type EnvironmentPut struct {
	InfrastructureRepoURL string                `json:"infrastructureRepoURL,omitempty"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules,omitempty"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN,omitempty"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

// EnvironmentPost is the implementation of the TowerAPI EnvironmentPostBody schema
type EnvironmentPost struct {
	Branch                string                `json:"branch,omitempty"`
	InfrastructureRepoURL string                `json:"infrastructureRepoURL,omitempty"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules,omitempty"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules,omitempty"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN,omitempty"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

// EnvironmentStatus is the implementation of the TowerAPI EnvironmentStatus schema
type EnvironmentStatus struct {
	Repository string `json:"repository,omitempty"`
	Branch     string `json:"branch,omitempty"`
	Status     string `json:"status,omitempty"`
}

// GitHubWebhook struct contains the three important values for auto-staging from the GitHub Webhook.
//
// ref_type must be branch for auto-staging
//
// ref is the name of the Git Branch
//
// repository/name is the name of the repository
type GitHubWebhook struct {
	Ref        string `json:"ref"`
	RefType    string `json:"ref_type"`
	Repository struct {
		Name string `json:"name"`
	}
}

// TriggerSchedulePost is the implementation of the TowerAPI EnvironmentStatus schema
type TriggerSchedulePost struct {
	Branch     string `json:"branch"`
	Repository string `json:"repository"`
	Action     string `json:"action"`
}

// TimeSchedule is the implementation of the TowerAPI TimeSchedule schema
type TimeSchedule struct {
	Cron string `json:"cron"`
}

// InternalServerErrorResponse contains a APIGatewayProxyResponse struct preset with "Internal server error" it's used as return value in controllers.
var InternalServerErrorResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Internal server error\"}",
	StatusCode: 500,
}

// InvalidRequestBodyResponse contains a APIGatewayProxyResponse struct preset with "Invalid request body" it's used as return value in controllers.
var InvalidRequestBodyResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Invalid request body\"}",
	StatusCode: 400,
}

// InvalidWebhookIsDeactivatedResponse contains a APIGatewayProxyResponse struct preset with "Webhooks are deactivated for this repository" it's used as return value in controllers.
var InvalidWebhookIsDeactivatedResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Webhooks are deactivated for this repository\"}",
	StatusCode: 400,
}

// InvalidEnvironmentStatusResponse contains a APIGatewayProxyResponse struct preset with "Can't execute operation in current environment status" it's used as return value in controllers.
var InvalidEnvironmentStatusResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Can't execute operation in current environment status\"}",
	StatusCode: 400,
}

// NotFoundErrorResponse contains a APIGatewayProxyResponse struct preset with "Not found" it's used as return value in controllers.
var NotFoundErrorResponse = events.APIGatewayProxyResponse{
	Body:       "{\"message\": \"Not found\"}",
	StatusCode: 404,
}
