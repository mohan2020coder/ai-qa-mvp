package workflow

import (
	"database/sql"
	"encoding/json"
	"log"

	"agentic-sdlc-api/internal/agents"
)

func LoadWorkflowFromDB(db *sql.DB, llm agents.LLMInterface, workflowID string) (*Workflow, error) {
	// Load workflow metadata first (optional: name/description)
	var wfNameNS sql.NullString
	err := db.QueryRow("SELECT name FROM workflows WHERE id=$1", workflowID).Scan(&wfNameNS)
	if err != nil {
		return nil, err
	}
	wfName := ""
	if wfNameNS.Valid {
		wfName = wfNameNS.String
	}

	rows, err := db.Query(`
		SELECT id, name, type, config, next_ids, condition 
		FROM workflow_nodes 
		WHERE workflow_id = $1
	`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make(map[string]*Node)

	for rows.Next() {
		var (
			id, nodeType string
			nodeNameNS   sql.NullString
			conditionNS  sql.NullString
			configJSON   []byte
			nextIDsJSON  []byte
		)

		if err := rows.Scan(&id, &nodeNameNS, &nodeType, &configJSON, &nextIDsJSON, &conditionNS); err != nil {
			log.Println("row scan error:", err)
			continue
		}

		nodeName := ""
		if nodeNameNS.Valid {
			nodeName = nodeNameNS.String
		}

		condition := ""
		if conditionNS.Valid {
			condition = conditionNS.String
		}

		var config map[string]any
		if len(configJSON) > 0 {
			if err := json.Unmarshal(configJSON, &config); err != nil {
				log.Println("config unmarshal error:", err)
			}
		} else {
			config = make(map[string]any)
		}

		var nextIDs []string
		if len(nextIDsJSON) > 0 {
			if err := json.Unmarshal(nextIDsJSON, &nextIDs); err != nil {
				log.Println("next_ids unmarshal error:", err)
			}
		}

		// Create agent based on node type
		var agent agents.AgentInterface
		switch nodeType {
		case "llm":
			systemPrompt, _ := config["system"].(string)
			agent = agents.NewLLMAgent(nodeName, llm, systemPrompt, 2048)

		case "api":
			agent = agents.NewAPIAgent(nodeName, config)

		case "script":
			agent = agents.NewScriptAgent(nodeName, config)

		default:
			log.Println("unknown node type:", nodeType)
			continue
		}

		nodes[id] = &Node{
			ID:        id,
			Agent:     agent,
			Config:    config,
			Next:      nextIDs,
			Condition: condition,
		}
	}

	return &Workflow{
		ID:    workflowID,
		Name:  wfName,
		Nodes: nodes,
	}, nil
}
