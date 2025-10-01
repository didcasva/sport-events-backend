package repository

import (
	"time"

	"sport-events-backend/internal/config"
	"sport-events-backend/internal/models"
)

// CreateEvent inserta un evento y devuelve su id.
func CreateEvent(e models.Event) (int, error) {
	var id int
	query := `
		INSERT INTO events (name, description, type, date, location, route, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := config.DB.QueryRow(query, e.Name, e.Description, e.Type, e.Date, e.Location, e.Route, e.CreatedBy).Scan(&id)
	return id, err
}

// GetAllEvents obtiene todos los eventos como []models.Event
func GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	query := `SELECT id, name, description, type, date, location, route, created_by, created_at FROM events`
	err := config.DB.Select(&events, query)
	return events, err
}

// RegisterUserToEvent registra a un usuario en un evento.
func RegisterUserToEvent(userID, eventID int) error {
	query := `INSERT INTO registrations (user_id, event_id, date) VALUES ($1, $2, $3)`
	_, err := config.DB.Exec(query, userID, eventID, time.Now())
	return err
}
// GetRegistrationsByEvent obtiene todas las inscripciones para un evento específico.
func GetRegistrationsByEvent(eventID int) ([]models.Registration, error) {
	var regs []models.Registration
	// Primero obtenemos las inscripciones básicas
	query := `SELECT id, user_id, event_id, date FROM registrations WHERE event_id = $1`
	if err := config.DB.Select(&regs, query, eventID); err != nil {
		return nil, err
	}

	// Para cada registro, cargamos la info del usuario asociado
	for i := range regs {
		var u models.User
		if err := config.DB.Get(&u, `SELECT id, name, email, role, created_at FROM users WHERE id = $1`, regs[i].UserID); err != nil {
			return nil, err
		}
		regs[i].User = u
	}

	return regs, nil
}


