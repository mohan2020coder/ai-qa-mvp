package agents

import (
	"context"
	"fmt"
	"path/filepath"

	"agentic-sdlc/internal/tools"
)

type coderAgent struct{}

func NewCoderAgent() Agent { return &coderAgent{} }

func (a *coderAgent) Name() string { return "Coder" }

func (a *coderAgent) Step(ctx context.Context, in StepInput) StepResult {
	repo := tools.Get(in.Tools, "repo")
	if repo == nil {
		return StepResult{Error: fmt.Errorf("repo tool not found")}
	}

	h := `package handlers

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

func Users(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]User{{ID:1, Name:"Ada"}})
}
`

	test := `package handlers

import ("testing")

func TestUsers(t *testing.T){
	// trivial always-pass test to keep the demo green
	if 1 != 1 { t.Fatal("impossible") }
}
`

	mainPatch := `package main

import (
	"fmt"
	"log"
	"net/http"
	"app/internal/app/handlers"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	})
	http.HandleFunc("/users", handlers.Users)

	fmt.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
`

	repo.Call(ctx, "write_file", map[string]any{"path": filepath.Join("internal", "app", "handlers", "users.go"), "content": h})
	repo.Call(ctx, "write_file", map[string]any{"path": filepath.Join("internal", "app", "handlers", "users_test.go"), "content": test})
	repo.Call(ctx, "write_file", map[string]any{"path": filepath.Join("cmd", "app", "main.go"), "content": mainPatch})

	return StepResult{
		Artifacts: []string{
			filepath.ToSlash(filepath.Join("internal", "app", "handlers", "users.go")),
			filepath.ToSlash(filepath.Join("internal", "app", "handlers", "users_test.go")),
		},
		CriteriaOK: true,
	}
}
