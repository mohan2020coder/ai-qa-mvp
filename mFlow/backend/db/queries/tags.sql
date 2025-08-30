-- name: ListTags
SELECT * FROM tags;

-- name: GetTagByID
SELECT * FROM tags WHERE id = $1;

-- name: CreateTag
INSERT INTO tags (name)
VALUES ($1)
RETURNING *;

-- name: DeleteTag
DELETE FROM tags WHERE id = $1;
