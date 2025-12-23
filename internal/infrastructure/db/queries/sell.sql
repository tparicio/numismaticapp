-- name: MarkCoinAsSold :one
UPDATE coins
SET
    sold_at = $2,
    sold_price = $3,
    sale_channel = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetDistinctSaleChannels :many
SELECT DISTINCT sale_channel
FROM coins
WHERE sale_channel IS NOT NULL AND sale_channel != ''
ORDER BY sale_channel;
