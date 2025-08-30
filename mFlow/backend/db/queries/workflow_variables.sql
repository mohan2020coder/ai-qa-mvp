-- name: ListWorkflowVariablesByWorkflow
SELECT * FROM workflow_variables WHERE workflow_id = $1;

-- name: GetWorkflowVariableByID
SELECT * FROM workflow_variables WHERE id = $1;

-- name: CreateWorkflowVariable
INSERT INTO workflow_variables (workflow_id, key, value, is_secret)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateWorkflowVariable
UPDATE workflow_variables
SET value = $2, is_secret = $3
WHERE id = $1
RETURNING *;

-- name: DeleteWorkflowVariable
DELETE FROM workflow_variables WHERE id = $1;
