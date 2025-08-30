-- name: ListNodesByWorkflow
SELECT * FROM nodes WHERE workflow_id = $1;

-- name: GetNodeByID
SELECT * FROM nodes WHERE id = $1;

-- name: CreateNode
INSERT INTO nodes (workflow_id, type, name, position_x, position_y, parameters, credentials_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateNode
UPDATE nodes
SET parameters = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteNode
DELETE FROM nodes WHERE id = $1;
