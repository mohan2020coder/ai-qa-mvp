package agents

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"agentic-sdlc/internal/tools"
)

type architectAgent struct{}

func NewArchitectAgent() Agent { return &architectAgent{} }

func (a *architectAgent) Name() string { return "Architect" }

func (a *architectAgent) Step(ctx context.Context, in StepInput) StepResult {
	_ = ctx
	repo := tools.Get(in.Tools, "repo")
	if repo == nil {
		return StepResult{Error: fmt.Errorf("repo tool not found")}
	}

	adr := `# ADR-0001: Service Architecture

## Context
A minimal HTTP service with one /health and one /users endpoint.

## Decision
- Go 1.22
- Standard library HTTP (router can be added later)
- Docker image built via Dockerfile
- Optional Postgres via env var DATABASE_URL

## Consequences
- Simple, fast startup
- Easy to extend with middlewares
`

	openapi := `openapi: 3.0.0
info:
  title: Sample App
  version: 0.1.0
paths:
  /health:
    get:
      responses:
        '200':
          description: ok
  /users:
    get:
      responses:
        '200':
          description: list users
`

	repo.Call(ctx, "write_file", map[string]any{"path": filepath.Join("docs", "ADR-0001.md"), "content": adr})
	repo.Call(ctx, "write_file", map[string]any{"path": filepath.Join("api", "openapi.yaml"), "content": openapi})

	return StepResult{
		Artifacts: []string{
			filepath.ToSlash(filepath.Join("docs", "ADR-0001.md")),
			filepath.ToSlash(filepath.Join("api", "openapi.yaml")),
		},
		CriteriaOK: strings.Contains(adr, "Go 1.22"),
	}
}
