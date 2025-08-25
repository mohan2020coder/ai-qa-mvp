package agents

import (
	"agentic-sdlc-ollama/internal/contracts"
	"github.com/tmc/langchaingo/llms"
)

func NewProductAgent(llm llms.Model) contracts.Agent {
	return &baseLLMAgent{
		name: "01-Product",
		llm:  llm,
		system: `You are a senior Product Manager AI.
Write clear, concise *product requirements* from the given spec.
Deliver: problem statement, personas, top 5 user stories (As a … I want … so that …), acceptance criteria, and non-functional requirements.
Keep it under 400 lines.`,
	}
}
