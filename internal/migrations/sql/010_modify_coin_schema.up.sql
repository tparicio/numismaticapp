-- Rename columns
ALTER TABLE coins RENAME COLUMN notes TO technical_notes;
ALTER TABLE coins RENAME COLUMN user_notes TO personal_notes;

-- Change mintage type
ALTER TABLE coins ALTER COLUMN mintage TYPE BIGINT;

-- Add physical properties
ALTER TABLE coins ADD COLUMN weight_g NUMERIC;
ALTER TABLE coins ADD COLUMN diameter_mm NUMERIC;
ALTER TABLE coins ADD COLUMN thickness_mm NUMERIC;
ALTER TABLE coins ADD COLUMN edge TEXT;
ALTER TABLE coins ADD COLUMN shape TEXT;

-- Add transaction data
ALTER TABLE coins ADD COLUMN acquired_at DATE;
ALTER TABLE coins ADD COLUMN sold_at DATE;
ALTER TABLE coins ADD COLUMN price_paid NUMERIC(10, 2);
ALTER TABLE coins ADD COLUMN sold_price NUMERIC(10, 2);

-- Drop sample image columns
ALTER TABLE coins DROP COLUMN sample_image_url_front;
ALTER TABLE coins DROP COLUMN sample_image_url_back;
