package repository

import "sport-events-backend/internal/config"

func CountRegistrationsForEvent(eventID int) (int, error) {
	const q = `SELECT COUNT(*) FROM registrations WHERE event_id = $1`
	var total int
	err := config.DB.Get(&total, q, eventID)
	return total, err
}

func CancelRegistration(userID, eventID int) error {
    query := `
        DELETE FROM registrations
        WHERE user_id = $1 AND event_id = $2
    `
    _, err := config.DB.Exec(query, userID, eventID)
    return err
}

