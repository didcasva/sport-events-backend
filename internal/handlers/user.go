package handlers

import (
	"encoding/json"
	"net/http"

	"sport-events-backend/internal/middleware"
	
	"sport-events-backend/internal/repository"
)

func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Perfil base
	user, err := repository.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Respuesta rol-aware
	resp := map[string]interface{}{
		"profile": user,
	}

	switch user.Role {
	case "organizer":
		events, err := repository.GetEventsByCreator(user.ID)
		if err == nil {
			resp["events_created"] = events
		} else {
			resp["events_created"] = []interface{}{}
		}
	default: // runner / athlete / cualquier otro
		regs, err := repository.GetUserRegistrationsWithEvents(user.ID)
		if err == nil {
			resp["my_registrations"] = regs
		} else {
			resp["my_registrations"] = []interface{}{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}