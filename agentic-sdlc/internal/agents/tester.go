package agents

import (
	"context"
	"fmt"

	"agentic-sdlc/internal/tools"
)

type testerAgent struct{}

func NewTesterAgent() Agent { return &testerAgent{} }

func (a *testerAgent) Name() string { return "Tester" }

func (a *testerAgent) Step(ctx context.Context, in StepInput) StepResult {
	testTool := tools.Get(in.Tools, "test")
	if testTool == nil {
		return StepResult{Error: fmt.Errorf("test tool not found")}
	}
	res, err := testTool.Call(ctx, "go_test", map[string]any{})
	if err != nil {
		return StepResult{Error: err}
	}
	ok := false
	if passed, _ := res.Data["passed"].(bool); passed {
		ok = true
	}
	return StepResult{Artifacts: res.Artifacts, CriteriaOK: ok}
}
