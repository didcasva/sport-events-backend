package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
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


func RegisterEventHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
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

	err = repository.RegisterUserToEvent(claims.UserID, eventID)
	if err != nil {
		// Duplicado por UNIQUE (user_id, event_id)
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			http.Error(w, "Ya estabas inscrito en este evento", http.StatusConflict) // 409
			return
		}
		http.Error(w, "Error registrando usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Usuario inscrito correctamente",
	})
}
func CancelRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	okDel, err := repository.CancelUserRegistration(claims.UserID, eventID)
	if err != nil {
		http.Error(w, "Error cancelando inscripción: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !okDel {
		http.Error(w, "No estabas inscrito en este evento", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Inscripción cancelada",
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
func GetMyRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := repository.GetUserRegistrationsWithEvents(claims.UserID)
	if err != nil {
		http.Error(w, "Error obteniendo mis inscripciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}
// PUT /api/events/{id}  (solo organizer dueño)
func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	// Parse input completo (PUT espera todos los campos)
	var in struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Type        string          `json:"type"`
		Date        time.Time       `json:"date"`
		Location    string          `json:"location"`
		Route       json.RawMessage `json:"route"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validaciones mínimas
	if in.Name == "" || in.Type == "" || in.Location == "" {
		http.Error(w, "Faltan campos obligatorios (name, type, location)", http.StatusBadRequest)
		return
	}
	// (opcional) no permitir fecha en el pasado
	// if in.Date.Before(time.Now().Add(-time.Minute)) { ... }

	e := models.Event{
		ID:          eventID,
		Name:        in.Name,
		Description: in.Description,
		Type:        in.Type,
		Date:        in.Date,
		Location:    in.Location,
		Route:       in.Route,
	}

	okUpd, err := repository.UpdateEventByOwner(e, claims.UserID)
	if err != nil || !okUpd {
		http.Error(w, "No autorizado para editar o evento inexistente", http.StatusForbidden)
		return
	}

	updated, _ := repository.GetEventByID(eventID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DELETE /api/events/{id}  (solo organizer dueño)
func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	okDel, err := repository.DeleteEventByOwner(eventID, claims.UserID)
	if err != nil || !okDel {
		http.Error(w, "No autorizado para eliminar o evento inexistente", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Evento eliminado"})
}

func GetEventDetailHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	evt, err := repository.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Evento no encontrado", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"event": evt,
	}

	// Si es organizer y dueño del evento, incluir inscritos
	if claims.Role == "organizer" && evt.CreatedBy == claims.UserID {
		regs, err := repository.GetRegistrationsForEvent(eventID)
		if err == nil {
			resp["registrations"] = regs
		} else {
			resp["registrations"] = []interface{}{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
// GET /api/events?type=&location=&date= consultas por filtros
func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	eventType := r.URL.Query().Get("type")
	location := r.URL.Query().Get("location")
	date := r.URL.Query().Get("date")

	var events []models.Event
	var err error

	if eventType != "" || location != "" || date != "" {
		events, err = repository.GetEventsFiltered(eventType, location, date)
	} else {
		events, err = repository.GetAllEvents()
	}

	if err != nil {
		http.Error(w, "Error obteniendo eventos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
