package at

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/chiyoi/apricot/kitsune"
	"github.com/chiyoi/iter/res"
	"github.com/gorilla/schema"
)

type OptionsThreadsCreate struct {
	Messages []MessagePost
	Metadata map[string]string
}

func ThreadsCreate(ctx context.Context, o *OptionsThreadsCreate) (thread Thread, err error) {
	u, err := url.JoinPath(Endpoint, "threads")
	token, err := res.R(ctx, err, getToken)
	var body io.Reader
	if o != nil {
		var req any = struct {
			Messages []MessagePost     `json:"messages,omitempty"`
			Metadata map[string]string `json:"metadata,omitempty"`
		}{}
		body, err = res.R(req, err, kitsune.JSONReader)
	}
	err = res.C(u, err, kitsune.Post(ctx, body, &thread, setHeader(token)))
	return
}

func ThreadsMessagesList(ctx context.Context, threadID string, query QueryPager) (pager Pager[Message], err error) {
	u, err := url.JoinPath(Endpoint, "threads", threadID, "messages")
	token, err := res.R(ctx, err, getToken)
	err = res.C(u, err, kitsune.Get(ctx, &pager, res.ComposedHooks(
		setHeader(token),
		setQuery(query),
	)))
	return
}

type OptionsThreadsMessagesCreate struct {
	FileIDs  []string
	Metadata map[string]string
}

func ThreadsMessagesCreate(ctx context.Context, threadID string, role Role, content string, o *OptionsThreadsMessagesCreate) (message Message, err error) {
	var req struct {
		Role     Role              `json:"role"`
		Content  string            `json:"content"`
		FileIDs  []string          `json:"file_ids,omitempty"`
		Metadata map[string]string `json:"metadata,omitempty"`
	}
	req.Role = role
	req.Content = content
	if o != nil {
		req.FileIDs = o.FileIDs
		req.Metadata = o.Metadata
	}
	u, err := url.JoinPath(Endpoint, "threads", threadID, "messages")
	token, err := res.R(ctx, err, getToken)
	body, err := res.R(any(req), err, kitsune.JSONReader)
	err = res.C(u, err, kitsune.Post(ctx, body, &message, res.ComposedHooks(
		kitsune.SetHeaderContentTypeJSON,
		setHeader(token),
	)))
	return
}

type OptionsThreadsRunsCreate struct {
	Model        string
	Instructions string
	Tools        []map[string]any
	Metadata     map[string]string
}

func ThreadsRunsCreate(ctx context.Context, threadID, assistantID string, o *OptionsThreadsRunsCreate) (run Run, err error) {
	var req struct {
		AssistantID  string            `json:"assistant_id"`
		Model        string            `json:"model,omitempty"`
		Instructions string            `json:"instructions,omitempty"`
		Tools        []map[string]any  `json:"tools,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
	}
	req.AssistantID = assistantID
	if o != nil {
		req.Model = o.Model
		req.Instructions = o.Instructions
		req.Tools = o.Tools
		req.Metadata = o.Metadata
	}
	u, err := url.JoinPath(Endpoint, "threads", threadID, "runs")
	token, err := res.R(ctx, err, getToken)
	body, err := res.R(any(req), err, kitsune.JSONReader)
	err = res.C(u, err, kitsune.Post(ctx, body, &run, res.ComposedHooks(
		kitsune.SetHeaderContentTypeJSON,
		setHeader(token),
	)))
	return
}

func ThreadsRunsRetrieve(ctx context.Context, threadID string, runID string) (run Run, err error) {
	u, err := url.JoinPath(Endpoint, "threads", threadID, "runs", runID)
	token, err := res.R(ctx, err, getToken)
	err = res.C(u, err, kitsune.Get(ctx, &run, setHeader(token)))
	return
}

func ThreadsRunsCancel(ctx context.Context, threadID, runID string) (run Run, err error) {
	u, err := url.JoinPath(Endpoint, "threads", threadID, "runs", runID, "cancel")
	token, err := res.R(ctx, err, getToken)
	err = res.C(u, err, kitsune.Post(ctx, nil, &run, setHeader(token)))
	return
}

func ThreadsRunsStepsList(ctx context.Context, threadID string, runID string, query QueryPager) (pager Pager[RunStep], err error) {
	u, err := url.JoinPath(Endpoint, "threads", threadID, "runs", runID, "steps")
	token, err := res.R(ctx, err, getToken)
	err = res.C(u, err, kitsune.Get(ctx, &pager, res.ComposedHooks(
		setHeader(token),
		setQuery(query),
	)))
	return
}

func ThreadsRunsToolOutputsSubmit(ctx context.Context, threadID, runID string, outputs []ToolOutput) (run Run, err error) {
	req := struct {
		ToolOutputs []ToolOutput `json:"tool_outputs"`
	}{outputs}
	u, err := url.JoinPath(Endpoint, "threads", threadID, "runs", runID, "submit_tool_outputs")
	token, err := res.R(ctx, err, getToken)
	body, err := res.R(any(req), err, kitsune.JSONReader)
	err = res.C(u, err, kitsune.Post(ctx, body, &run, res.ComposedHooks(
		kitsune.SetHeaderContentTypeJSON,
		setHeader(token),
	)))
	return
}

func getToken(ctx context.Context) (token string, err error) {
	token, ok := ctx.Value(ContextKeyAPIKey).(string)
	if !ok {
		err = errors.New("missing context (" + string(ContextKeyAPIKey) + ")")
	}
	return
}

func setHeader(token string) func(r *http.Request) (*http.Request, error) {
	return func(r *http.Request) (*http.Request, error) {
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("OpenAI-Beta", "assistants=v1")
		return r, nil
	}
}

func setQuery(query any) func(r *http.Request) (*http.Request, error) {
	return func(r *http.Request) (*http.Request, error) {
		vs := make(url.Values)
		schema.NewEncoder().Encode(query, vs)
		q := r.URL.Query()
		for k, v := range vs {
			q[k] = v
		}
		r.URL.RawQuery = q.Encode()
		return r, nil
	}
}

type Pager[Data any] struct {
	Data    []Data `json:"data"`
	FirstID string `json:"first_id"`
	LastID  string `json:"last_id"`
	HasMore bool   `json:"has_more"`
}

type QueryPager struct {
	Limit  int    `schema:"limit,omitempty"`
	Order  string `schema:"order,omitempty"`
	After  string `schema:"after,omitempty"`
	Before string `schema:"before,omitempty"`
}
