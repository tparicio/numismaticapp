DROP TYPE IF EXISTS image_type CASCADE;
CREATE TYPE image_type AS ENUM ('original_front', 'original_back', 'processed_front', 'processed_back');

CREATE TABLE coin_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coin_id UUID NOT NULL REFERENCES coins(id) ON DELETE CASCADE,
    image_type image_type NOT NULL,
    path VARCHAR(255) NOT NULL,
    extension VARCHAR(10) NOT NULL,
    size BIGINT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    mime_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coin_images_coin_id ON coin_images(coin_id);
