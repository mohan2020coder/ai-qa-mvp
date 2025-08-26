package agents

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

type AgentInterface interface {
	Name() string
	Run(ctx context.Context, input string) (string, error)
	MaxTokens() int
	SetMaxTokens(n int)
}

// in agents/interface.go
type LLMInterface = llms.Model
