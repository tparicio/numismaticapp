-- name: CreateCoin :one
INSERT INTO coins (
    id, name, mint, mintage, country, year, face_value, currency, material, description, km_code,
    min_value, max_value, grade, technical_notes, gemini_details, group_id, personal_notes,
    weight_g, diameter_mm, thickness_mm, edge, shape,
    acquired_at, sold_at, price_paid, sold_price, numista_number, numista_details,
    gemini_model, gemini_temperature, numista_search,
    ruler, orientation, series, commemorated_topic
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16, $17, $18,
    $19, $20, $21, $22, $23,
    $24, $25, $26, $27, $28, $29,
    $30, $31, $32,
    $33, $34, $35, $36
) RETURNING *;

-- name: GetCoin :one
SELECT * FROM coins
WHERE id = $1 LIMIT 1;

-- name: ListCoins :many
SELECT * FROM coins
WHERE 
    (sqlc.narg('group_id')::int IS NULL OR group_id = sqlc.narg('group_id'))
    AND (sqlc.narg('year')::int IS NULL OR year = sqlc.narg('year'))
    AND (sqlc.narg('country')::text IS NULL OR country ILIKE sqlc.narg('country'))
    AND (sqlc.narg('query')::text IS NULL OR 
        name ILIKE '%' || sqlc.narg('query') || '%' OR 
        description ILIKE '%' || sqlc.narg('query') || '%' OR
        km_code ILIKE '%' || sqlc.narg('query') || '%'
    )
    AND (sqlc.narg('min_price')::float8 IS NULL OR min_value >= sqlc.narg('min_price')::float8)
    AND (sqlc.narg('max_price')::float8 IS NULL OR max_value <= sqlc.narg('max_price')::float8)
ORDER BY
    CASE WHEN sqlc.narg('sort_by')::text = 'year' AND sqlc.narg('sort_order')::text = 'asc' THEN year END ASC,
    CASE WHEN sqlc.narg('sort_by')::text = 'year' AND sqlc.narg('sort_order')::text = 'desc' THEN year END DESC,
    
    CASE WHEN sqlc.narg('sort_by')::text = 'min_value' AND sqlc.narg('sort_order')::text = 'asc' THEN min_value END ASC,
    CASE WHEN sqlc.narg('sort_by')::text = 'min_value' AND sqlc.narg('sort_order')::text = 'desc' THEN min_value END DESC,
    
    CASE WHEN sqlc.narg('sort_by')::text = 'max_value' AND sqlc.narg('sort_order')::text = 'asc' THEN max_value END ASC,
    CASE WHEN sqlc.narg('sort_by')::text = 'max_value' AND sqlc.narg('sort_order')::text = 'desc' THEN max_value END DESC,
    
    CASE WHEN sqlc.narg('sort_by')::text = 'created_at' AND sqlc.narg('sort_order')::text = 'asc' THEN created_at END ASC,
    CASE WHEN sqlc.narg('sort_by')::text = 'created_at' AND sqlc.narg('sort_order')::text = 'desc' THEN created_at END DESC,
    
    CASE WHEN sqlc.narg('sort_by')::text = 'country' AND sqlc.narg('sort_order')::text = 'asc' THEN country END ASC,
    CASE WHEN sqlc.narg('sort_by')::text = 'country' AND sqlc.narg('sort_order')::text = 'desc' THEN country END DESC,
    
    CASE WHEN sqlc.narg('sort_by')::text = 'name' AND sqlc.narg('sort_order')::text = 'asc' THEN name END ASC,
    CASE WHEN sqlc.narg('sort_by')::text = 'name' AND sqlc.narg('sort_order')::text = 'desc' THEN name END DESC,
    
    created_at DESC
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
    numista_number = $28,
    numista_details = $29,
    gemini_model = $30,
    gemini_temperature = $31,
    numista_search = $32,
    ruler = $33,
    orientation = $34,
    series = $35,
    commemorated_topic = $36,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCoin :exec
DELETE FROM coins
WHERE id = $1;

-- name: GetTotalValue :one
SELECT COALESCE(SUM(max_value), 0)::float8 FROM coins;

-- name: GetAverageValue :one
SELECT COALESCE(AVG(max_value), 0)::float8 FROM coins;

-- name: ListTopValuableCoins :many
SELECT * FROM coins
ORDER BY max_value DESC
LIMIT 5;

-- name: ListRecentCoins :many
SELECT * FROM coins
ORDER BY created_at DESC
LIMIT 5;

-- name: GetMaterialDistribution :many
SELECT material, COUNT(*) as count
FROM coins
WHERE material IS NOT NULL AND material != ''
GROUP BY material
ORDER BY count DESC;

-- name: GetGradeDistribution :many
SELECT grade, COUNT(*) as count
FROM coins
WHERE grade IS NOT NULL
GROUP BY grade
ORDER BY count DESC;

-- name: GetAllValues :many
SELECT max_value FROM coins WHERE max_value IS NOT NULL;

-- name: GetCountryDistribution :many
SELECT country, COUNT(*) as count FROM coins GROUP BY country;

-- name: GetOldestCoin :one
SELECT * FROM coins WHERE year > 0 ORDER BY year ASC LIMIT 1;

-- name: GetRarestCoins :many
SELECT * FROM coins WHERE mintage > 0 ORDER BY mintage ASC LIMIT $1;

-- name: GetGroupDistribution :many
SELECT COALESCE(g.name, 'Uncategorized') as group_name, COUNT(c.id) as count 
FROM coins c 
LEFT JOIN groups g ON c.group_id = g.id 
GROUP BY g.name;

-- name: GetTotalWeightByMaterial :one
SELECT COALESCE(SUM(weight_g), 0)::float8 FROM coins WHERE material ILIKE $1;

-- name: GetHeaviestCoin :one
SELECT * FROM coins ORDER BY weight_g DESC LIMIT 1;

-- name: GetSmallestCoin :one
SELECT * FROM coins WHERE diameter_mm > 0 ORDER BY diameter_mm ASC LIMIT 1;

-- name: GetRandomCoin :one
SELECT * FROM coins ORDER BY RANDOM() LIMIT 1;

-- name: GetAllCoins :many
SELECT * FROM coins;

-- name: GetGroupStats :many
SELECT 
    g.id as group_id, 
    COALESCE(g.name, 'Uncategorized') as group_name, 
    COUNT(c.id) as count,
    COALESCE(MIN(c.min_value), 0)::float8 as min_val,
    COALESCE(MAX(c.max_value), 0)::float8 as max_val
FROM coins c 
LEFT JOIN groups g ON c.group_id = g.id 
GROUP BY g.id, g.name
ORDER BY count DESC;
