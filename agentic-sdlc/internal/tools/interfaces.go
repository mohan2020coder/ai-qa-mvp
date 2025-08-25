package tools

import (
	"context"
	"fmt"
)

type CallResult struct {
	OK        bool
	Content   string
	Data      map[string]any
	Artifacts []string
}

type Tool interface {
	Name() string
	Description() string
	Call(ctx context.Context, input string, args map[string]any) (CallResult, error)
}

func Get(set []Tool, name string) Tool {
	for _, t := range set {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

func ok(content string, data map[string]any, artifacts []string) (CallResult, error) {
	return CallResult{OK: true, Content: content, Data: data, Artifacts: artifacts}, nil
}
func fail(msg string) (CallResult, error) {
	return CallResult{OK: false, Content: msg, Data: map[string]any{}, Artifacts: nil}, fmt.Errorf(msg)
}
