-- name: ListConnectionsByWorkflow
SELECT * FROM connections WHERE workflow_id = $1;

-- name: CreateConnection
INSERT INTO connections (workflow_id, from_node_id, to_node_id, output_index, input_index)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteConnection
DELETE FROM connections WHERE id = $1;
