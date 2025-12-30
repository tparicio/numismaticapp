-- name: GetCoinPercentiles :one
WITH stats AS (
    SELECT 
        c.id,
        c.max_value,
        c.mintage,
        c.weight_g,
        c.diameter_mm,
        PERCENT_RANK() OVER (ORDER BY c.max_value) as value_percentile,
        PERCENT_RANK() OVER (ORDER BY c.mintage DESC) as rarity_percentile, -- Lower mintage = higher rarity, so DESC for "rarity score"
        PERCENT_RANK() OVER (ORDER BY c.weight_g) as weight_percentile,
        PERCENT_RANK() OVER (ORDER BY c.diameter_mm) as size_percentile
    FROM coins c
    WHERE c.sold_at IS NULL -- Compare only with collection
)
SELECT 
    value_percentile,
    rarity_percentile,
    weight_percentile,
    size_percentile
FROM stats
WHERE id = $1;

-- name: GetCollectionYearDistribution :many
SELECT 
    year, 
    COUNT(*) as count
FROM coins
WHERE 
    year > 0 
    AND sold_at IS NULL
GROUP BY year
ORDER BY year;

-- name: GetCollectionGradeDistribution :many
SELECT 
    grade, 
    COUNT(*) as count
FROM coins
WHERE 
    grade != '' 
    AND sold_at IS NULL
GROUP BY grade
ORDER BY grade;
