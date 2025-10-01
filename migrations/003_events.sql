CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  name VARCHAR(150) NOT NULL,
  description TEXT,
  type VARCHAR(50),
  date TIMESTAMP NOT NULL,
  location VARCHAR(150),
  route JSONB,
  created_by INT REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW()
);
