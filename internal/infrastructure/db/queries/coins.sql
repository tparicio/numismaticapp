-- name: CreateCoin :one
INSERT INTO coins (
    country, year, face_value, currency, material, description, km_code,
    min_value, max_value, grade, sample_image_url_front, sample_image_url_back,
    notes, gemini_details, group_id, user_notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
) RETURNING *;

-- name: GetCoin :one
SELECT * FROM coins
WHERE id = $1 LIMIT 1;

-- name: ListCoins :many
SELECT * FROM coins
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountCoins :one
SELECT count(*) FROM coins;

-- name: UpdateCoin :one
UPDATE coins
SET
    country = $2,
    year = $3,
    face_value = $4,
    currency = $5,
    material = $6,
    description = $7,
    km_code = $8,
    min_value = $9,
    max_value = $10,
    grade = $11,
    sample_image_url_front = $12,
    sample_image_url_back = $13,
    notes = $14,
    gemini_details = $15,
    group_id = $16,
    user_notes = $17,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCoin :exec
DELETE FROM coins
WHERE id = $1;
