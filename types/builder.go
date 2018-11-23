package types

type BuilderEvent struct {
	Repository            string                `json:"repository"`
	Branch                string                `json:"branch"`
	Operation             string                `json:"operation"`
	InfrastructureRepoURL string                `json:"infrastructureRepoUrl"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
	Success               int                   `json:"success"`
}
