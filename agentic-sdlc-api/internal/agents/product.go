package agents

import "github.com/tmc/langchaingo/llms"

func NewProductAgent(llm llms.Model) AgentInterface {
	return &BaseLLMAgent{
		name: "01-Product",
		llm:  llm,
		system: `You are a senior Product Manager AI.
        Write clear, concise *product requirements* from the given spec.
        Deliver: problem statement, personas, top 5 user stories (As a … I want … so that …), 
        acceptance criteria, and non-functional requirements.
        Keep it under 1000 lines.`,
	}
}
