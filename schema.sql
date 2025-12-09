CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE coins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,
    face_value VARCHAR(100) NOT NULL,
    currency VARCHAR(100) NOT NULL,
    material VARCHAR(100) NOT NULL,
    description TEXT,
    km_code VARCHAR(50),
    min_value DECIMAL(10, 2),
    max_value DECIMAL(10, 2),
    grade VARCHAR(50),
    sample_image_url_front TEXT,
    sample_image_url_back TEXT,
    notes TEXT,
    gemini_details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coins_country ON coins(country);
CREATE INDEX idx_coins_year ON coins(year);
