-- name: ListCredentialsByOwner
SELECT * FROM credentials WHERE owner_id = $1;

-- name: GetCredentialByID
SELECT * FROM credentials WHERE id = $1;

-- name: CreateCredential
INSERT INTO credentials (name, type, data, iv, tag, owner_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteCredential
DELETE FROM credentials WHERE id = $1;
