package agents

import (
	"context"
	"fmt"

	"agentic-sdlc/internal/tools"
)

type reviewerAgent struct{}

func NewReviewerAgent() Agent { return &reviewerAgent{} }

func (a *reviewerAgent) Name() string { return "Reviewer" }

func (a *reviewerAgent) Step(ctx context.Context, in StepInput) StepResult {
	lint := tools.Get(in.Tools, "lint")
	if lint == nil {
		return StepResult{Error: fmt.Errorf("lint tool not found")}
	}
	res, err := lint.Call(ctx, "lint", map[string]any{})
	if err != nil {
		return StepResult{Error: err}
	}
	return StepResult{Artifacts: res.Artifacts, CriteriaOK: true}
}
