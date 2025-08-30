-- name: ListWorkflowVersionsByWorkflow
SELECT * FROM workflow_versions WHERE workflow_id = $1 ORDER BY version DESC;

-- name: GetWorkflowVersionByID
SELECT * FROM workflow_versions WHERE id = $1;

-- name: CreateWorkflowVersion
INSERT INTO workflow_versions (workflow_id, version, definition)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteWorkflowVersion
DELETE FROM workflow_versions WHERE id = $1;
