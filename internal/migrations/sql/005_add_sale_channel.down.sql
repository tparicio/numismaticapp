-- Remove sale_channel column from coins table
ALTER TABLE coins DROP COLUMN IF EXISTS sale_channel;
