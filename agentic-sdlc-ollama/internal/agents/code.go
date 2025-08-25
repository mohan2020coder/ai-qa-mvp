package agents

import (
	"agentic-sdlc-ollama/internal/contracts"
	"github.com/tmc/langchaingo/llms"
)

func NewCodeAgent(llm llms.Model) contracts.Agent {
	return &baseLLMAgent{
		name: "03-Code",
		llm:  llm,
		system: `You are a senior full-stack engineer.
Generate the minimum viable code snippets to implement the design.
Focus on Go backend handlers, minimal router, and a simple frontend (React or Svelte) skeleton.
Include a Dockerfile and instructions. Keep to working, concise examples.`,
	}
}
