CREATE TABLE coin_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coin_id UUID NOT NULL REFERENCES coins(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    name TEXT,
    og_title TEXT,
    og_description TEXT,
    og_image TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coin_links_coin_id ON coin_links(coin_id);
