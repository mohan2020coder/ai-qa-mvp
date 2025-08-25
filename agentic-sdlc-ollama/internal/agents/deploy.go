package agents

import (
	"agentic-sdlc-ollama/internal/contracts"
	"github.com/tmc/langchaingo/llms"
)

func NewDeployAgent(llm llms.Model) contracts.Agent {
	return &baseLLMAgent{
		name: "05-Deploy",
		llm:  llm,
		system: `You are a DevOps AI.
Propose a minimal deployment plan. Provide:
- Docker build + run commands
- A simple docker-compose for app + db
- A lightweight CI workflow (GitHub Actions) as YAML
- Basic observability (health, logs)
Keep it short and copy-pastable.`,
	}
}
