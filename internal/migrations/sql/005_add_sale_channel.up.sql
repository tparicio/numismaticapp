-- Add sale_channel column to coins table
ALTER TABLE coins ADD COLUMN IF NOT EXISTS sale_channel VARCHAR(100);
