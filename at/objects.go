package at

import "strings"

type Thread struct {
	ID       string         `json:"id"`
	Object   string         `json:"object"`
	CreateAt int64          `json:"create_at"`
	Metadata map[string]any `json:"metadata"`
}

type Message struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int    `json:"created_at"`
	ThreadID  string `json:"thread_id"`
	Role      Role   `json:"role"`
	Content   []struct {
		Type MessageContentElementType `json:"type"`
		Text struct {
			Value       string `json:"value"`
			Annotations []any  `json:"annotations"`
		} `json:"text"`
	} `json:"content"`
	FileIds     []string       `json:"file_ids"`
	AssistantID string         `json:"assistant_id"`
	RunID       string         `json:"run_id"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

func (message *Message) Plaintext() string {
	var buf strings.Builder
	for _, paragraph := range message.Content {
		if paragraph.Type == MessageContentElementTypeText {
			buf.WriteString(paragraph.Text.Value)
		}
	}
	return buf.String()
}

type MessageContentElementType string

const (
	MessageContentElementTypeImageFile MessageContentElementType = "image_file"
	MessageContentElementTypeText      MessageContentElementType = "text"
)

type MessagePost struct {
	// Role only support `user` currently.
	Role     Role           `json:"role"`
	Content  string         `json:"content"`
	FileIds  []string       `json:"file_ids"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Run struct {
	ID           string           `json:"id"`
	Object       string           `json:"object"`
	CreatedAt    int64            `json:"created_at"`
	AssistantID  string           `json:"assistant_id"`
	ThreadID     string           `json:"thread_id"`
	Status       RunStatus        `json:"status"`
	StartedAt    int64            `json:"started_at"`
	ExpiresAt    int64            `json:"expires_at"`
	CancelledAt  int64            `json:"cancelled_at"`
	FailedAt     int64            `json:"failed_at"`
	CompletedAt  int64            `json:"completed_at"`
	LastError    *Error           `json:"last_error"`
	Model        string           `json:"model"`
	Instructions string           `json:"instructions"`
	Tools        []map[string]any `json:"tools"`
	FileIds      []string         `json:"file_ids"`
	Metadata     map[string]any   `json:"metadata"`
}

type RunStep struct {
	ID          string         `json:"id"`
	Object      string         `json:"object"`
	CreatedAt   int64          `json:"created_at"`
	AssistantID string         `json:"assistant_id"`
	ThreadID    string         `json:"thread_id"`
	RunID       string         `json:"run_id"`
	Type        RunStepType    `json:"type"`
	Status      RunStatus      `json:"status"`
	CancelledAt int64          `json:"cancelled_at"`
	CompletedAt int64          `json:"completed_at"`
	ExpiredAt   int64          `json:"expired_at"`
	FailedAt    int64          `json:"failed_at"`
	LastError   *Error         `json:"last_error"`
	StepDetails StepDetails    `json:"step_details"`
	Metadata    map[string]any `json:"metadata"`
}

type StepDetails struct {
	Type            string          `json:"type"`
	MessageCreation MessageCreation `json:"message_creation"`
	ToolCalls       []ToolCall      `json:"tool_calls"`
}

type MessageCreation struct {
	MessageID string `json:"message_id"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function,omitempty"`
	// TODO: Other call types
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Parameters struct {
	Properties any      `json:"properties"`
	Required   []string `json:"required"`
}
