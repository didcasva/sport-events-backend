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
	Status            string          `db:"status" json:"status"`
	CancelledAt       *time.Time      `db:"cancelled_at" json:"cancelled_at,omitempty"`
	CancellationReason *string        `db:"cancellation_reason" json:"cancellation_reason,omitempty"`

}
type EventSummary struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	Date      time.Time `db:"date" json:"date"`
	Location  string    `db:"location" json:"location"`
	CreatedBy int       `db:"created_by" json:"created_by"`
}

type RegistrationWithEvent struct {
	RegistrationID int       `db:"registration_id" json:"registration_id"`
	EventID        int       `db:"event_id" json:"event_id"`
	Name           string    `db:"name" json:"name"`
	Type           string    `db:"type" json:"type"`
	Date           time.Time `db:"date" json:"date"`
	Location       string    `db:"location" json:"location"`
	RegisteredAt   time.Time `db:"registered_at" json:"registered_at"`
}

type EventRegistrationUser struct {
	RegistrationID int    `db:"registration_id" json:"registration_id"`
	UserID         int    `db:"user_id" json:"user_id"`
	UserName       string `db:"user_name" json:"user_name"`
	UserEmail      string `db:"user_email" json:"user_email"`
}
