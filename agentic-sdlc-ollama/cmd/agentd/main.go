package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"agentic-sdlc-ollama/internal/agents"
	"agentic-sdlc-ollama/internal/contracts"
	"agentic-sdlc-ollama/internal/orchestrator"

	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	var specPath string
	var model string
	var out string
	flag.StringVar(&specPath, "spec", "examples/sample-spec.md", "Path to the product spec file")
	flag.StringVar(&model, "model", getenv("OLLAMA_MODEL", "llama3"), "Ollama model name (e.g. llama3, codellama, mistral)")
	flag.StringVar(&out, "out", ".workspace", "Output directory for artifacts")
	flag.Parse()

	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		log.Fatalf("read spec: %v", err)
	}
	spec := strings.TrimSpace(string(specBytes))

	// Connect to Ollama
	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		log.Fatalf("ollama: %v", err)
	}

	// Wire agents
	var list []contracts.Agent
	list = append(list, agents.NewProductAgent(llm))
	list = append(list, agents.NewDesignAgent(llm))
	list = append(list, agents.NewCodeAgent(llm))
	list = append(list, agents.NewTestAgent(llm))
	list = append(list, agents.NewDeployAgent(llm))

	orc := orchestrator.New(list, out)

	fmt.Printf("== Agentic SDLC with Ollama (%s) ==\n", model)
	ctx := context.Background()
	report, err := orc.Execute(ctx, spec)
	if err != nil {
		log.Fatalf("execute: %v", err)
	}

	fmt.Println("\nArtifacts:")
	for _, r := range report.Results {
		fmt.Printf(" - %s → %s\n", r.Agent, r.Path)
	}
	fmt.Printf("Done. Combined output in %s\n", out+"/outputs/combined.md")
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
