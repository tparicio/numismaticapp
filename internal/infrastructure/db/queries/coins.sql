-- name: CreateCoin :one
INSERT INTO coins (
    id, name, mint, mintage, country, year, face_value, currency, material, description, km_code,
    min_value, max_value, grade, technical_notes, gemini_details, group_id, personal_notes,
    weight_g, diameter_mm, thickness_mm, edge, shape,
    acquired_at, sold_at, price_paid, sold_price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16, $17, $18,
    $19, $20, $21, $22, $23,
    $24, $25, $26, $27
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
    technical_notes = $15,
    gemini_details = $16,
    group_id = $17,
    personal_notes = $18,
    weight_g = $19,
    diameter_mm = $20,
    thickness_mm = $21,
    edge = $22,
    shape = $23,
    acquired_at = $24,
    sold_at = $25,
    price_paid = $26,
    sold_price = $27,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCoin :exec
DELETE FROM coins
WHERE id = $1;
