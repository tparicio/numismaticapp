-- name: CreateCoin :one
INSERT INTO coins (
    id, name, mint, mintage, country, year, face_value, currency, material, description, km_code,
    min_value, max_value, grade, sample_image_url_front, sample_image_url_back,
    notes, gemini_details, group_id, user_notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16, $17, $18, $19, $20
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
    name = $2,
    mint = $3,
    mintage = $4,
    country = $5,
    year = $6,
    face_value = $7,
    currency = $8,
    material = $9,
    description = $10,
    km_code = $11,
    min_value = $12,
    max_value = $13,
    grade = $14,
    sample_image_url_front = $15,
    sample_image_url_back = $16,
    notes = $17,
    gemini_details = $18,
    group_id = $19,
    user_notes = $20,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCoin :exec
DELETE FROM coins
WHERE id = $1;
