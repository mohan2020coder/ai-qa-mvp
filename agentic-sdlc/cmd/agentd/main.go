package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"agentic-sdlc/internal/orchestrator"
	"agentic-sdlc/internal/agents"
	"agentic-sdlc/internal/tools"
)

func main() {
	ctx := context.Background()
	ws := ".workspace"
	_ = os.MkdirAll(ws, 0o755)

	specBytes, err := os.ReadFile("examples/sample-spec.md")
	if err != nil { log.Fatalf("read spec: %v", err) }
	spec := string(specBytes)

	repo := tools.NewRepoTool(filepath.Join(ws, "repo"))
	testTool := tools.NewTestTool(filepath.Join(ws, "repo"))
	lintTool := tools.NewLintTool(filepath.Join(ws, "repo"))
	buildTool := tools.NewBuildTool(filepath.Join(ws, "repo"), filepath.Join(ws, "artifacts"))
	deployTool := tools.NewDeployTool(filepath.Join(ws, "artifacts"))
	toolset := []tools.Tool{repo, testTool, lintTool, buildTool, deployTool}

	reg := agents.Registry{
		Architect:  agents.NewArchitectAgent(),
		Scaffolder: agents.NewScaffolderAgent(),
		Coder:      agents.NewCoderAgent(),
		Tester:     agents.NewTesterAgent(),
		Reviewer:   agents.NewReviewerAgent(),
		DevOps:     agents.NewDevOpsAgent(),
	}
	planner := orchestrator.NewSimplePlanner()

	fmt.Println("== Agentic SDLC: Plan → Execute ==")
	dag := planner.Plan(ctx, spec)
	fmt.Printf("Planned %d tasks\n", len(dag.Nodes))

	start := time.Now()
	res := orchestrator.Execute(ctx, dag, reg, toolset, ws)
	fmt.Printf("\nCompleted in %s\n", time.Since(start))
	fmt.Printf("Artifacts:\n")
	for _, a := range res.Artifacts {
		fmt.Printf(" - %s\n", a)
	}
}
