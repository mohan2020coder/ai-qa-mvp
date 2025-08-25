package agents

import (
	"context"
	"fmt"
	"path/filepath"

	"agentic-sdlc/internal/tools"
)

type scaffolderAgent struct{}

func NewScaffolderAgent() Agent { return &scaffolderAgent{} }

func (a *scaffolderAgent) Name() string { return "Scaffolder" }

func (a *scaffolderAgent) Step(ctx context.Context, in StepInput) StepResult {
	repo := tools.Get(in.Tools, "repo")
	if repo == nil {
		return StepResult{Error: fmt.Errorf("repo tool not found")}
	}

	goMod := "module app\n\ngo 1.22\n"
	mainGo := `package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	})
	fmt.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
`

	repo.Call(ctx, "write_file", map[string]any{"path": "go.mod", "content": goMod})
	repo.Call(ctx, "write_file", map[string]any{"path": filepath.Join("cmd", "app", "main.go"), "content": mainGo})
	repo.Call(ctx, "write_file", map[string]any{"path": "Dockerfile", "content": "FROM golang:1.22-alpine\nWORKDIR /app\nCOPY . .\nRUN go build -o app ./cmd/app\nCMD [\"./app\"]\n"})

	return StepResult{
		Artifacts: []string{
			"go.mod",
			filepath.ToSlash(filepath.Join("cmd", "app", "main.go")),
			"Dockerfile",
		},
		CriteriaOK: true,
	}
}
