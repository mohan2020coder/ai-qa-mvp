-- name: ListExecutionsByWorkflow
SELECT * FROM executions WHERE workflow_id = $1;

-- name: GetExecutionByID
SELECT * FROM executions WHERE id = $1;

-- name: CreateExecution
INSERT INTO executions (workflow_id, status, execution_data, run_by_user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateExecutionStatus
UPDATE executions
SET status = $2, ended_at = now(), error_message = $3
WHERE id = $1
RETURNING *;
