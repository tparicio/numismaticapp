-- name: CreateCoinImage :one
INSERT INTO coin_images (
    coin_id,
    image_type,
    side,
    path,
    extension,
    size,
    width,
    height,
    mime_type,
    original_filename
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: ListCoinImagesByCoinID :many
SELECT * FROM coin_images
WHERE coin_id = $1
ORDER BY created_at ASC;

-- name: ListCoinImagesByCoinIDs :many
SELECT * FROM coin_images
WHERE coin_id = ANY($1::uuid[])
ORDER BY coin_id, created_at ASC;
