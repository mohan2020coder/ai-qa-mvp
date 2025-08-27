package outdated

// import (
// 	"log"

// 	"github.com/tmc/langchaingo/llms/ollama"
// )

// func BuildAgents(llm *ollama.LLM, names []string) []AgentInterface {
// 	var agentsList []AgentInterface
// 	for _, name := range names {
// 		switch name {
// 		case "Product":
// 			agentsList = append(agentsList, NewProductAgent(llm))
// 		case "Design":
// 			agentsList = append(agentsList, NewDesignAgent(llm))
// 		case "Code":
// 			agentsList = append(agentsList, NewCodeAgent(llm))
// 		case "Test":
// 			agentsList = append(agentsList, NewTestAgent(llm))
// 		case "Deploy":
// 			agentsList = append(agentsList, NewDeployAgent(llm))
// 		default:
// 			log.Printf("⚠️ Unknown agent name: %s (skipping)", name)
// 		}
// 	}
// 	return agentsList
// }
