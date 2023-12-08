package at

type RunStatus string

const (
	// Only for Runs
	RunStatusQueued         RunStatus = "queued"
	RunStatusRequiresAction RunStatus = "requires_action"
	RunStatusCancelling     RunStatus = "cancelling"

	// For Runs and RunSteps
	RunStatusInProgress RunStatus = "in_progress"
	RunStatusCancelled  RunStatus = "cancelled"
	RunStatusFailed     RunStatus = "failed"
	RunStatusCompleted  RunStatus = "completed"
	RunStatusExpired    RunStatus = "expired"
)

type Role string

const (
	RoleUser Role = "user"
)

type RunStepType string

const (
	RunStepTypeMessageCreation RunStepType = "message_creation"
	RunStepTypeToolCalls       RunStepType = "tool_calls"
)
