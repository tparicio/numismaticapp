CREATE TABLE group_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_group_images_group_id ON group_images(group_id);

CREATE TABLE coin_gallery_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coin_id UUID NOT NULL REFERENCES coins(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coin_gallery_images_coin_id ON coin_gallery_images(coin_id);
