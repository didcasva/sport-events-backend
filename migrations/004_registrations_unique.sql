-- migrations/004_registrations_unique.sql
ALTER TABLE registrations
ADD CONSTRAINT uniq_user_event UNIQUE (user_id, event_id);
