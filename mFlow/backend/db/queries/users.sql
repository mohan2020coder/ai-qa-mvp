-- name: ListUsers
SELECT * FROM users;

-- name: GetUserByID
SELECT * FROM users WHERE id = $1;

-- name: CreateUser
INSERT INTO users (email, password_hash, name, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser
UPDATE users
SET name = $2, role = $3, is_active = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser
DELETE FROM users WHERE id = $1;
