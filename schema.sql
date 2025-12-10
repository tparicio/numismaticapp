CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE image_type AS ENUM ('original', 'crop', 'thumbnail', 'sample');
CREATE TYPE coin_side AS ENUM ('front', 'back');
CREATE TYPE grade_type AS ENUM ('MC', 'RC', 'BC', 'MBC', 'EBC', 'SC', 'FDC', 'PROOF');

CREATE TABLE coins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255),
    mint VARCHAR(255),
    mintage INTEGER,
    country VARCHAR(255),
    year INTEGER,
    face_value VARCHAR(100),
    currency VARCHAR(100),
    material VARCHAR(100),
    description TEXT,
    km_code VARCHAR(50),
    min_value DECIMAL(10, 2),
    max_value DECIMAL(10, 2),
    grade grade_type,
    sample_image_url_front TEXT,
    sample_image_url_back TEXT,
    notes TEXT,
    gemini_details JSONB,
    group_id INT REFERENCES groups(id),
    user_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
