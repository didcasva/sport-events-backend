package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"github.com/gorilla/mux"

	"sport-events-backend/internal/middleware"
	"sport-events-backend/internal/models"
	"sport-events-backend/internal/repository"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Solo organizadores pueden crear eventos
	if claims.Role != "organizer" {
		http.Error(w, "Forbidden: role 'organizer' required", http.StatusForbidden)
		return
	}

	var input struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Type        string          `json:"type"`
		Date        time.Time       `json:"date"`
		Location    string          `json:"location"`
		Route       json.RawMessage `json:"route"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := models.Event{
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		Date:        input.Date,
		Location:    input.Location,
		Route:       input.Route,
		CreatedBy:   claims.UserID,
	}

	id, err := repository.CreateEvent(event)
	if err != nil {
		http.Error(w, "Error creando evento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}


func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := repository.GetAllEvents() // función que consultas la BD
	if err != nil {
		http.Error(w, "Error obteniendo eventos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func RegisterEventHandler(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	err = repository.RegisterUserToEvent(userClaims.UserID, eventID)
	if err != nil {
		http.Error(w, "Error registrando usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Usuario inscrito correctamente",
	})
}
// internal/handlers/events.go
func GetEventRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	registrations, err := repository.GetRegistrationsByEvent(eventID)
	if err != nil {
		http.Error(w, "Error obteniendo registros: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registrations)
}
