-- name: ListWebhooksByTrigger
SELECT * FROM webhooks WHERE trigger_id = $1;

-- name: GetWebhookByID
SELECT * FROM webhooks WHERE id = $1;

-- name: CreateWebhook
INSERT INTO webhooks (trigger_id, path, method, secret)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteWebhook
DELETE FROM webhooks WHERE id = $1;
