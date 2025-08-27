package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIAgent struct {
	name   string
	config map[string]any
	maxTok int
}

func NewAPIAgent(name string, config map[string]any) *APIAgent {
	return &APIAgent{name: name, config: config}
}

func (a *APIAgent) Name() string {
	return a.name
}

func (a *APIAgent) Run(ctx context.Context, input string) (string, error) {
	url, ok := a.config["url"].(string)
	if !ok || url == "" {
		return "", fmt.Errorf("missing API url")
	}

	method := "POST"
	if m, ok := a.config["method"].(string); ok {
		method = m
	}

	// Prepare payload
	payload := map[string]any{
		"input": input,
	}
	if extra, ok := a.config["payload"].(map[string]any); ok {
		for k, v := range extra {
			payload[k] = v
		}
	}

	bodyBytes, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil
}

func (a *APIAgent) MaxTokens() int {
	return a.maxTok
}

func (a *APIAgent) SetMaxTokens(n int) {
	a.maxTok = n
}
