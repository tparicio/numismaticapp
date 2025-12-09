DROP TYPE IF EXISTS grade_type CASCADE;
CREATE TYPE grade_type AS ENUM ('MC', 'RC', 'BC', 'MBC', 'EBC', 'SC', 'FDC', 'PROOF');

-- Since DROP TYPE CASCADE drops the column, we need to add it back.
ALTER TABLE coins ADD COLUMN grade grade_type;

COMMENT ON TYPE grade_type IS 'Standardized coin grades: MC, RC, BC, MBC, EBC, SC, FDC, PROOF';
