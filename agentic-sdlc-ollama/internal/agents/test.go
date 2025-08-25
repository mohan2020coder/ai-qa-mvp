package agents

import (
	"agentic-sdlc-ollama/internal/contracts"
	"github.com/tmc/langchaingo/llms"
)

func NewTestAgent(llm llms.Model) contracts.Agent {
	return &baseLLMAgent{
		name: "04-Test",
		llm:  llm,
		system: `You are a Test Engineer AI.
Create a compact test plan and example tests.
Include unit tests (Go), API contract tests (curl examples), and a brief CI outline.
Return runnable Go tests for handlers when possible.`,
	}
}
