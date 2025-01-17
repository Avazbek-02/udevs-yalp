CREATE TYPE notification_status AS ENUM ('read', 'unread');

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY UNIQUE NOT NULL,
    owner_id UUID,
    user_id UUID REFERENCES users(id),
    email TEXT,
    message TEXT,
    status notification_status DEFAULT 'unread',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
