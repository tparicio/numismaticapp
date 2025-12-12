DROP TYPE IF EXISTS grade_type CASCADE;
CREATE TYPE grade_type AS ENUM ('MC', 'RC', 'BC', 'MBC', 'EBC', 'SC', 'FDC', 'PROOF');

-- Since DROP TYPE CASCADE drops the column ONLY IF it uses the type, 
-- but if it's VARCHAR (from initial schema), it won't be dropped.
-- So we explicitlly drop it to be safe and ensure we can add it back as the ENUM type.
ALTER TABLE coins DROP COLUMN IF EXISTS grade;
ALTER TABLE coins ADD COLUMN grade grade_type;

COMMENT ON TYPE grade_type IS 'Standardized coin grades: MC, RC, BC, MBC, EBC, SC, FDC, PROOF';
