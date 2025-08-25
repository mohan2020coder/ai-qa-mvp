package agents

import (
	"context"
	"fmt"

	"agentic-sdlc/internal/tools"
)

type devopsAgent struct{}

func NewDevOpsAgent() Agent { return &devopsAgent{} }

func (a *devopsAgent) Name() string { return "DevOps" }

func (a *devopsAgent) Step(ctx context.Context, in StepInput) StepResult {
	build := tools.Get(in.Tools, "build")
	deploy := tools.Get(in.Tools, "deploy")
	if build == nil || deploy == nil {
		return StepResult{Error: fmt.Errorf("build/deploy tools not found")}
	}
	bres, err := build.Call(ctx, "docker_build", map[string]any{"name": "agentic/app", "tag": "dev"})
	if err != nil {
		return StepResult{Error: err}
	}
	dres, err := deploy.Call(ctx, "deploy_staging", map[string]any{"image": bres.Data["image"]})
	if err != nil {
		return StepResult{Error: err}
	}
	arts := append([]string{}, bres.Artifacts...)
	arts = append(arts, dres.Artifacts...)
	return StepResult{Artifacts: arts, CriteriaOK: true}
}
