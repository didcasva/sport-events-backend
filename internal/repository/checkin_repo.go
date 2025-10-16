package repository

import (
	"sport-events-backend/internal/config"
	"time"
)

type Checkin struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	EventID     int       `db:"event_id" json:"event_id"`
	CheckpointID int      `db:"checkpoint_id" json:"checkpoint_id"`
	Lat         float64   `db:"lat" json:"lat"`
	Lng         float64   `db:"lng" json:"lng"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func CreateCheckin(userID, eventID, checkpointID int, lat, lng float64) error {
	query := `
        INSERT INTO checkins (user_id, event_id, checkpoint_id, lat, lng)
        VALUES ($1, $2, $3, $4, $5)`
	_, err := config.DB.Exec(query, userID, eventID, checkpointID, lat, lng)
	return err
}
