-- Drop everything to ensure clean slate (User requested reset)
DROP TABLE IF EXISTS schema_migrations CASCADE;
DROP TABLE IF EXISTS coin_images CASCADE;
DROP TABLE IF EXISTS coins CASCADE;
DROP TABLE IF EXISTS groups CASCADE;
DROP TYPE IF EXISTS image_type CASCADE;
DROP TYPE IF EXISTS coin_side CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Ensure schema_migrations exists with the correct schema
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty BOOLEAN NOT NULL DEFAULT false,
    executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE image_type AS ENUM ('original', 'crop', 'thumbnail', 'sample');
CREATE TYPE coin_side AS ENUM ('front', 'back');

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE coins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255),
    mint VARCHAR(255),
    mintage BIGINT,
    country VARCHAR(255),
    year INTEGER,
    face_value VARCHAR(100),
    currency VARCHAR(100),
    material VARCHAR(100),
    description TEXT,
    km_code VARCHAR(50),
    min_value DECIMAL(10, 2),
    max_value DECIMAL(10, 2),
    grade VARCHAR(50),
    technical_notes TEXT,
    gemini_details JSONB,
    group_id INT REFERENCES groups(id),
    personal_notes TEXT,
    weight_g NUMERIC,
    diameter_mm NUMERIC,
    thickness_mm NUMERIC,
    edge TEXT,
    shape TEXT,
    numista_number INTEGER,
    acquired_at DATE,
    sold_at DATE,
    price_paid NUMERIC(10, 2),
    sold_price NUMERIC(10, 2),
    sale_channel VARCHAR(100),
    gemini_model VARCHAR(100),
    gemini_temperature NUMERIC(3, 2),
    
    -- Consolidated columns from previous migrations
    numista_search TEXT,   -- from 002
    ruler TEXT NOT NULL DEFAULT '', -- from 003
    orientation TEXT NOT NULL DEFAULT '', -- from 003
    series TEXT NOT NULL DEFAULT '', -- from 003
    commemorated_topic TEXT NOT NULL DEFAULT '', -- from 003
    numista_details JSONB, -- from 004

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coins_country ON coins(country);
CREATE INDEX idx_coins_year ON coins(year);

CREATE TABLE coin_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coin_id UUID NOT NULL REFERENCES coins(id) ON DELETE CASCADE,
    image_type image_type NOT NULL,
    side coin_side NOT NULL,
    path VARCHAR(255) NOT NULL,
    extension VARCHAR(10) NOT NULL,
    size BIGINT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    mime_type VARCHAR(50) NOT NULL,
    original_filename VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coin_images_coin_id ON coin_images(coin_id);
