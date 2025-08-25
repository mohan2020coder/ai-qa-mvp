package agents

import (
	"agentic-sdlc-ollama/internal/contracts"
	"github.com/tmc/langchaingo/llms"
)

func NewDesignAgent(llm llms.Model) contracts.Agent {
	return &baseLLMAgent{
		name: "02-Design",
		llm:  llm,
		system: `You are a pragmatic Software Architect AI.
Design a minimal architecture for the product requirements.
Deliver: C4-style context, components, sequence for key flows, API spec outline, data model, and tradeoffs.
Output in markdown with code blocks for API examples.`,
	}
}
