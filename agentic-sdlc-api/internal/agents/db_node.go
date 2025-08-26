package agents

import (
	"context"
)

// DBNode implements AgentInterface for any agent type
type DBNode struct {
	agent AgentInterface
}

func (d *DBNode) Name() string {
	return d.agent.Name()
}

func (d *DBNode) Run(ctx context.Context, input string) (string, error) {
	return d.agent.Run(ctx, input)
}

func (d *DBNode) MaxTokens() int {
	return d.agent.MaxTokens()
}

func (d *DBNode) SetMaxTokens(n int) {
	d.agent.SetMaxTokens(n)
}
