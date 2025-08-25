package orchestrator

import (
	"agentic-sdlc/internal/agents"
	"context"
	"strings"
)

type Planner interface {
	Plan(ctx context.Context, spec string) agents.TaskDAG
}

type SimplePlanner struct{}

func NewSimplePlanner() *SimplePlanner { return &SimplePlanner{} }

func (p *SimplePlanner) Plan(ctx context.Context, spec string) agents.TaskDAG {
	_ = ctx
	nodes := []agents.TaskNode{
		{ID: "arch", Type: "ARCHITECT", Summary: "Produce ADR and API design", Outputs: []string{"docs/ADR-0001.md", "api/openapi.yaml"}},
		{ID: "scaf", Type: "SCAFFOLD", Inputs: []string{"docs/ADR-0001.md"}, Summary: "Init repo layout", Outputs: []string{"go.mod", "cmd/app/main.go"}},
		{ID: "code", Type: "CODE", Inputs: []string{"cmd/app/main.go", "api/openapi.yaml"}, Summary: "Implement /users handler + tests", Outputs: []string{"internal/app/handlers/users.go"}},
		{ID: "test", Type: "TEST", Inputs: []string{"internal/app/handlers/users.go"}, Summary: "Unit tests & coverage", Criteria: []string{"tests pass", "coverage >= 50%"}},
		{ID: "review", Type: "REVIEW", Inputs: []string{"reports/test.json"}, Summary: "Lint & static checks", Criteria: []string{"lints clean"}},
		{ID: "devops", Type: "DEVOPS", Inputs: []string{"Dockerfile"}, Summary: "Build image & deploy", Outputs: []string{"artifacts/image.txt", "deploy/staging.txt"}},
	}
	if strings.Contains(strings.ToLower(spec), "postgres") {
		nodes = append(nodes, agents.TaskNode{ID: "db", Type: "CODE", Inputs: []string{"docs/ADR-0001.md"}, Summary: "Add Postgres client and healthcheck", Outputs: []string{"internal/app/db/db.go"}})
	}
	edges := map[string][]string{
		"arch":   {"scaf"},
		"scaf":   {"code"},
		"code":   {"test"},
		"test":   {"review"},
		"review": {"devops"},
	}
	return agents.TaskDAG{Nodes: nodes, Edges: edges}
}
