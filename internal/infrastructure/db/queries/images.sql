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
    mime_type
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListCoinImagesByCoinID :many
SELECT * FROM coin_images
WHERE coin_id = $1
ORDER BY created_at ASC;
