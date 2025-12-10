-- name: CreateGroup :one
INSERT INTO groups (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetGroupByName :one
SELECT * FROM groups
WHERE name = $1;

-- name: ListGroups :many
SELECT * FROM groups
ORDER BY name;

-- name: UpdateGroup :one
UPDATE groups
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;
