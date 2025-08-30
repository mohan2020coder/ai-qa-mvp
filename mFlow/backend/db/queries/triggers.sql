-- name: ListTriggersByWorkflow
SELECT * FROM triggers WHERE workflow_id = $1;

-- name: GetTriggerByID
SELECT * FROM triggers WHERE id = $1;

-- name: CreateTrigger
INSERT INTO triggers (workflow_id, type, config, is_enabled)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTrigger
UPDATE triggers
SET config = $2, is_enabled = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTrigger
DELETE FROM triggers WHERE id = $1;
