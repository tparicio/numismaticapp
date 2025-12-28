-- name: AddCoinLink :one
INSERT INTO coin_links (
    coin_id, url, name, og_title, og_description, og_image
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: ListCoinLinks :many
SELECT * FROM coin_links
WHERE coin_id = $1
ORDER BY created_at DESC;

-- name: DeleteCoinLink :exec
DELETE FROM coin_links
WHERE id = $1;
