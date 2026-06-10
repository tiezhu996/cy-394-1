CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  friend_key VARCHAR(80),
  created_at TIMESTAMPTZ DEFAULT now()
);

INSERT INTO users (name, friend_key)
VALUES ('演示用户', 'friends-demo')
ON CONFLICT DO NOTHING;
