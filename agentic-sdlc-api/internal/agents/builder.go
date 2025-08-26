package agents

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/tmc/langchaingo/llms"
)

func BuildAgentsFromDB(db *sql.DB, llm llms.Model) ([]AgentInterface, error) {
	rows, err := db.Query("SELECT name, type, config FROM agents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agentsList []AgentInterface
	for rows.Next() {
		var name, agentType string
		var configJSON []byte
		if err := rows.Scan(&name, &agentType, &configJSON); err != nil {
			log.Println("scan error:", err)
			continue
		}

		var config map[string]any
		if err := json.Unmarshal(configJSON, &config); err != nil {
			log.Println("invalid config JSON for", name)
			continue
		}

		var agent AgentInterface
		switch agentType {
		case "LLM":
			systemPrompt, _ := config["system"].(string)
			maxTok := 2048
			if t, ok := config["max_tokens"].(float64); ok {
				maxTok = int(t)
			}
			agent = &BaseLLMAgent{name: name, llm: llm, system: systemPrompt, maxTok: maxTok}
		case "API":
			agent = &APIAgent{name: name, config: config}
		case "Script":
			agent = &ScriptAgent{name: name, config: config}
		default:
			log.Println("Unknown agent type:", agentType, "skipping")
			continue
		}

		agentsList = append(agentsList, &DBNode{agent: agent})
	}

	return agentsList, nil
}
