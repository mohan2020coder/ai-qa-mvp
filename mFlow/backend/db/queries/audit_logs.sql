-- name: ListAuditLogs
SELECT * FROM audit_logs ORDER BY timestamp DESC;

-- name: ListAuditLogsByEntity
SELECT * FROM audit_logs
WHERE entity_type = $1 AND entity_id = $2
ORDER BY timestamp DESC;

-- name: CreateAuditLog
INSERT INTO audit_logs (user_id, action, entity_type, entity_id, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
