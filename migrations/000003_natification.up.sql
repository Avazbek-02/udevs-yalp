CREATE TYPE notification_status AS ENUM ('read', 'unread');

CREATE TABLE if not exists notifications (
  id UUID PRIMARY KEY UNIQUE NOT NULL,
  user_id UUID REFERENCES users(id),
  message TEXT,
  status notification_status DEFAULT 'unread',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
