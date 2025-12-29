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

-- name: CreateGroupImage :one
INSERT INTO group_images (group_id, path)
VALUES ($1, $2)
RETURNING *;

-- name: ListGroupImages :many
SELECT * FROM group_images
WHERE group_id = $1
ORDER BY created_at ASC;

-- name: DeleteGroupImage :exec
DELETE FROM group_images
WHERE id = $1;

-- name: CreateCoinGalleryImage :one
INSERT INTO coin_gallery_images (coin_id, path)
VALUES ($1, $2)
RETURNING *;

-- name: ListCoinGalleryImages :many
SELECT * FROM coin_gallery_images
WHERE coin_id = $1
ORDER BY created_at ASC;

-- name: DeleteCoinGalleryImage :exec
DELETE FROM coin_gallery_images
WHERE id = $1;
