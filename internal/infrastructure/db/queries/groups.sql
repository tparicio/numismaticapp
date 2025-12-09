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
