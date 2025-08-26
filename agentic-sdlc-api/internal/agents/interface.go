package agents

import "context"

type AgentInterface interface {
	Name() string
	Run(ctx context.Context, input string) (string, error)
	MaxTokens() int
    SetMaxTokens(n int)

}
