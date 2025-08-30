-- name: ListAITasksByAgent
SELECT * FROM ai_tasks WHERE agent_id = $1;

-- name: GetAITaskByID
SELECT * FROM ai_tasks WHERE id = $1;

-- name: CreateAITask
INSERT INTO ai_tasks (agent_id, workflow_execution_id, input, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateAITaskStatus
UPDATE ai_tasks
SET status = $2, output = $3, error_message = $4, finished_at = now()
WHERE id = $1
RETURNING *;
