package models

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID          int             `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description"`
	Type        string          `db:"type" json:"type"`
	Date        time.Time       `db:"date" json:"date"`
	Location    string          `db:"location" json:"location"`
	Route       json.RawMessage `db:"route" json:"route"` // JSONB
	CreatedBy   int             `db:"created_by" json:"created_by"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
}
