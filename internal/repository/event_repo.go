package repository

import (
	"time"
	"errors"
	"fmt"
	
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
	query := `
		SELECT id, name, description, type, date, location, route, created_by, created_at, status
		FROM events
		WHERE status <> 'cancelled'
		ORDER BY date ASC
	`
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

func GetEventsByCreator(userID int) ([]models.EventSummary, error) {
	var evts []models.EventSummary
	query := `
		SELECT id, name, type, date, location, created_by
		FROM events
		WHERE created_by = $1
		ORDER BY date DESC;
	`
	err := config.DB.Select(&evts, query, userID)
	return evts, err
}
func GetUserRegistrationsWithEvents(userID int) ([]models.RegistrationWithEvent, error) {
	var rows []models.RegistrationWithEvent
	query := `
		SELECT 
			r.id         AS registration_id,
			e.id         AS event_id,
			e.name       AS name,
			e.type       AS type,
			e.date       AS date,
			e.location   AS location,
			r.date       AS registered_at
		FROM registrations r
		JOIN events e ON e.id = r.event_id
		WHERE r.user_id = $1
		ORDER BY e.date DESC;
	`
	err := config.DB.Select(&rows, query, userID)
	return rows, err
}

// Cancela inscripción; retorna (bool) si eliminó algo
func CancelUserRegistration(userID, eventID int) (bool, error) {
	const q = `
		DELETE FROM registrations
		WHERE user_id = $1 AND event_id = $2
		RETURNING id
	`
	var id int
	if err := config.DB.Get(&id, q, userID, eventID); err != nil {
		// No rows → no estaba inscrito
		// sqlx.Get retorna error si no hay fila; lo traducimos a (false, nil)
		// pero distinguimos el caso "no rows"
		// Para no importar paquetes extras, hacemos un check simple:
		// Si quieres fino: usar errors.Is(err, sql.ErrNoRows)
		return false, nil
	}
	return true, nil
}
// Obtener detalle (útil si quieres validar cosas antes de actualizar)
func GetEventByID(id int) (models.Event, error) {
	var e models.Event
	const q = `
		SELECT id, 
		name, 
		description, 
		type, 
		date, 
		location, 
		route, 
		created_by, 
		created_at,
		status,
		cancelled_at,
		cancellation_reason
		FROM events
		WHERE id = $1
	`
	err := config.DB.Get(&e, q, id)
	return e, err
}

// Actualizar solo si el owner coincide (WHERE id=? AND created_by=?)
func UpdateEventByOwner(e models.Event, ownerID int) (bool, error) {
	const q = `
		UPDATE events
		SET name = $1,
		    description = $2,
		    type = $3,
		    date = $4,
		    location = $5,
		    route = $6
		WHERE id = $7 AND created_by = $8
		RETURNING id
	`
	var id int
	if err := config.DB.Get(&id, q,
		e.Name, e.Description, e.Type, e.Date, e.Location, e.Route,
		e.ID, ownerID,
	); err != nil {
		// no rows → no es owner o no existe
		return false, err
	}
	return true, nil
}

// Eliminar solo si el owner coincide
func DeleteEventByOwner(eventID, ownerID int) (bool, error) {
	const q = `
		DELETE FROM events
		WHERE id = $1 AND created_by = $2
		RETURNING id
	`
	var id int
	if err := config.DB.Get(&id, q, eventID, ownerID); err != nil {
		return false, err
	}
	return true, nil
}

// (Opcional) Validación simple de existencia del evento
func MustOwnEvent(eventID, ownerID int) error {
	const q = `SELECT 1 FROM events WHERE id=$1 AND created_by=$2`
	var one int
	if err := config.DB.Get(&one, q, eventID, ownerID); err != nil {
		return errors.New("no eres dueño del evento o no existe")
	}
	return nil
}

func GetRegistrationsForEvent(eventID int) ([]models.EventRegistrationUser, error) {
	var rows []models.EventRegistrationUser
	const q = `
		SELECT 
			r.id   AS registration_id,
			u.id   AS user_id,
			u.name AS user_name,
			u.email AS user_email
		FROM registrations r
		JOIN users u ON u.id = r.user_id
		WHERE r.event_id = $1
		ORDER BY r.id DESC;
	`
	err := config.DB.Select(&rows, q, eventID)
	return rows, err
}


// Filtros opcionales: type, location (parcial), date (YYYY-MM-DD)
func GetEventsFiltered(eventType, location, date string,includeCancelled bool) ([]models.Event, error) {
	var events []models.Event
	query := `
		SELECT id, name, description, type, date, location, route, created_by, created_at,status
		FROM events
		WHERE 1=1
	`
	args := []interface{}{}
	i := 1

	if eventType != "" {
		query += fmt.Sprintf(" AND type = $%d", i)
		args = append(args, eventType)
		i++
	}
	if location != "" {
		query += fmt.Sprintf(" AND location ILIKE $%d", i)
		args = append(args, "%"+location+"%")
		i++
	}
	if date != "" {
		query += fmt.Sprintf(" AND DATE(date) = $%d", i)
		args = append(args, date) // formato YYYY-MM-DD
		i++
	}
	if !includeCancelled {
		query += " AND status <> 'cancelled'"
	}


	query += " ORDER BY date ASC"

	err := config.DB.Select(&events, query, args...)
	return events, err
}

// Cambiar estado a 'cancelled' solo si es owner
func CancelEventByOwner(eventID, ownerID int, reason string) (bool, error) {
	const q = `
		UPDATE events
		SET status = 'cancelled',
		    cancelled_at = NOW(),
		    cancellation_reason = $1
		WHERE id = $2 AND created_by = $3 AND status <> 'cancelled'
		RETURNING id
	`
	var id int
	if err := config.DB.Get(&id, q, reason, eventID, ownerID); err != nil {
		return false, err
	}
	return true, nil
}

// (útil para validaciones)
func GetEventStatus(eventID int) (string, error) {
	const q = `SELECT status FROM events WHERE id = $1`
	var st string
	err := config.DB.Get(&st, q, eventID)
	return st, err
}

