-- name: ListAIAgents
SELECT * FROM ai_agents;

-- name: GetAIAgentByID
SELECT * FROM ai_agents WHERE id = $1;

-- name: CreateAIAgent
INSERT INTO ai_agents (name, model, description, owner_id, config)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateAIAgent
UPDATE ai_agents
SET name = $2, description = $3, config = $4, is_active = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAIAgent
DELETE FROM ai_agents WHERE id = $1;
