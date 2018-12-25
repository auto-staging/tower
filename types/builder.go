package types

// BuilderEvent struct contains all values for the different Builder invoke bodys combined into one struct.
type BuilderEvent struct {
	Repository            string                `json:"repository"`
	Branch                string                `json:"branch"`
	Operation             string                `json:"operation"`
	InfrastructureRepoURL string                `json:"infrastructureRepoUrl"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
	Success               int                   `json:"success"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules"`
}
