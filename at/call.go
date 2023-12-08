package at

import (
	"context"
)

type FunctionCall struct {
	ctx             context.Context
	ThreadID, RunID string
	ToolCall        ToolCall
}

func (c *FunctionCall) Context() context.Context {
	return c.ctx
}

func NewFunctionCall(ctx context.Context, threadID, runID string, toolCall ToolCall) *FunctionCall {
	return &FunctionCall{
		ctx:      ctx,
		RunID:    runID,
		ThreadID: threadID,
		ToolCall: toolCall,
	}
}

type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}
