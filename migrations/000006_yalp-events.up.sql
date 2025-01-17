CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY not null,
    business_id UUID NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
    name VARCHAR(255),
    description TEXT,
    date text,
    location TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    UNIQUE(event_id, user_id)
);