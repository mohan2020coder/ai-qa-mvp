package agents

import (
	"context"

	"agentic-sdlc/internal/tools"
)

type StepInput struct {
	Node      TaskNode
	Tools     []tools.Tool
	Workspace string
}

type StepResult struct {
	Artifacts  []string
	CriteriaOK bool
	Error      error
}

type Agent interface {
	Name() string
	Step(ctx context.Context, in StepInput) StepResult
}

type Registry struct {
	Architect  Agent
	Scaffolder Agent
	Coder      Agent
	Tester     Agent
	Reviewer   Agent
	DevOps     Agent
}
type TaskNode struct {
	ID       string
	Type     string
	Inputs   []string
	Outputs  []string
	Criteria []string
	Summary  string
}

type TaskDAG struct {
	Nodes []TaskNode
	Edges map[string][]string
}
