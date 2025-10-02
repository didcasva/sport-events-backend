package handlers

import (
	"errors"
	"time"
)

func validateEventRequired(name, typ, location string, date time.Time) error {
	if name == "" || typ == "" || location == "" {
		return errors.New("faltan campos obligatorios (name, type, location)")
	}
	if date.IsZero() {
		return errors.New("date es obligatorio")
	}
	return nil
}

func validateEventFuture(date time.Time) error {
	// margen de 1 minuto para tolerar desfaces de reloj
	if date.Before(time.Now().Add(-1 * time.Minute)) {
		return errors.New("la fecha del evento debe ser futura")
	}
	return nil
}
