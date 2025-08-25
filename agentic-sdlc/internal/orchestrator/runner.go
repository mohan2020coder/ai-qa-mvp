package orchestrator

import (
	"context"
	"fmt"
	"sort"

	"agentic-sdlc/internal/agents"
	"agentic-sdlc/internal/tools"
)

type Result struct {
	Artifacts []string
}

func topoOrder(dag agents.TaskDAG) []agents.TaskNode {
	out := make([]agents.TaskNode, len(dag.Nodes))
	copy(out, dag.Nodes)
	return out
}

func Execute(ctx context.Context, dag agents.TaskDAG, reg agents.Registry, toolset []tools.Tool, workspace string) Result {
	order := topoOrder(dag)
	var artifacts []string

	for _, n := range order {
		fmt.Printf("\n→ [%s] %s\n", n.Type, n.Summary)
		var agent agents.Agent
		switch n.Type {
		case "ARCHITECT":
			agent = reg.Architect
		case "SCAFFOLD":
			agent = reg.Scaffolder
		case "CODE":
			agent = reg.Coder
		case "TEST":
			agent = reg.Tester
		case "REVIEW":
			agent = reg.Reviewer
		case "DEVOPS":
			agent = reg.DevOps
		default:
			fmt.Println("  (skip unknown task type)")
			continue
		}

		res := agent.Step(ctx, agents.StepInput{
			Node:      n,
			Tools:     toolset,
			Workspace: workspace,
		})
		if res.Error != nil {
			fmt.Printf("  ✖ error: %v\n", res.Error)
			continue
		}
		for _, a := range res.Artifacts {
			fmt.Printf("  ✓ artifact: %s\n", a)
			artifacts = append(artifacts, a)
		}
		if len(n.Criteria) > 0 {
			fmt.Printf("  criteria: %v → %v\n", n.Criteria, res.CriteriaOK)
		}
	}

	sort.Strings(artifacts)
	return Result{Artifacts: artifacts}
}
