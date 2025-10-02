-- migrations/005_events_status.sql
ALTER TABLE events
  ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'scheduled', -- scheduled | cancelled | completed
  ADD COLUMN cancelled_at TIMESTAMP NULL,
  ADD COLUMN cancellation_reason TEXT NULL;

-- Opcional: índice para filtrar rápido
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
