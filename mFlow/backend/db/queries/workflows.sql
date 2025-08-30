-- name: ListWorkflows
SELECT * FROM workflows;

-- name: GetWorkflowByID
SELECT * FROM workflows WHERE id = $1;

-- name: CreateWorkflow
INSERT INTO workflows (name, description, created_by, settings, is_template, forked_from)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateWorkflow
UPDATE workflows
SET name = $2, description = $3, settings = $4, is_active = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteWorkflow
DELETE FROM workflows WHERE id = $1;
