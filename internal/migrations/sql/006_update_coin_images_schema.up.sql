-- Create the new enum type for coin side
CREATE TYPE coin_side AS ENUM ('front', 'back');

-- Add the side column to coin_images
ALTER TABLE coin_images ADD COLUMN side coin_side;

-- Migrate existing data based on the old image_type values
UPDATE coin_images SET side = 'front' WHERE image_type::text LIKE '%_front';
UPDATE coin_images SET side = 'back' WHERE image_type::text LIKE '%_back';

-- Make side not null after migration
ALTER TABLE coin_images ALTER COLUMN side SET NOT NULL;

-- Rename the old image_type type to avoid conflict during transition (optional, but cleaner to just alter the column type)
-- However, PostgreSQL doesn't support ALTER TYPE ... ADD VALUE inside a transaction block easily if we want to completely replace the values.
-- Strategy:
-- 1. Create new enum type for image_type
-- 2. Add temporary column with new type
-- 3. Migrate data
-- 4. Drop old column
-- 5. Rename new column

CREATE TYPE new_image_type AS ENUM ('original', 'crop', 'thumbnail', 'sample');

ALTER TABLE coin_images ADD COLUMN new_image_type new_image_type;

UPDATE coin_images SET new_image_type = 'original' WHERE image_type::text LIKE 'original_%';
UPDATE coin_images SET new_image_type = 'crop' WHERE image_type::text LIKE 'processed_%';
-- Handle any other cases if necessary, or set default. Assuming only original and processed existed.

ALTER TABLE coin_images ALTER COLUMN new_image_type SET NOT NULL;

-- Drop the old column and type
ALTER TABLE coin_images DROP COLUMN image_type;
DROP TYPE image_type;

-- Rename new column to image_type
ALTER TABLE coin_images RENAME COLUMN new_image_type TO image_type;

-- Rename the new type to image_type (optional, but good for consistency)
ALTER TYPE new_image_type RENAME TO image_type;
