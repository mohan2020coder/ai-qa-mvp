-- name: ListWorkflowTags
SELECT wt.workflow_id, t.id AS tag_id, t.name
FROM workflow_tags wt
JOIN tags t ON t.id = wt.tag_id
WHERE wt.workflow_id = $1;

-- name: AddWorkflowTag
INSERT INTO workflow_tags (workflow_id, tag_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveWorkflowTag
DELETE FROM workflow_tags
WHERE workflow_id = $1 AND tag_id = $2;
