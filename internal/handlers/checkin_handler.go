package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"sport-events-backend/internal/middleware"
	"sport-events-backend/internal/repository"
)

const checkpointRadius = 3.0 // metros

// Haversine formula para distancia en metros
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

type Checkpoint struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Type string  `json:"type"`
}

func CheckinHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventID, _ := strconv.Atoi(vars["id"])
	checkpointID, _ := strconv.Atoi(vars["checkpointId"])

	var input struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Formato inválido", http.StatusBadRequest)
		return
	}

	// Obtener la ruta del evento
	route, err := repository.GetEventRoute(eventID)
	if err != nil {
		http.Error(w, "No se encontró la ruta", http.StatusNotFound)
		return
	}

	// Parsear JSON de la ruta con checkpoints
	var routeData struct {
		Checkpoints []Checkpoint `json:"checkpoints"`
	}
	if err := json.Unmarshal(route, &routeData); err != nil {
		http.Error(w, "Error parseando checkpoints", http.StatusInternalServerError)
		return
	}

	// Buscar el checkpoint
	var cp *Checkpoint
	for _, c := range routeData.Checkpoints {
		if c.ID == checkpointID {
			// copiar en variable local
			cp = &Checkpoint{
				ID:   c.ID,
				Name: c.Name,
				Lat:  c.Lat,
				Lng:  c.Lng,
				Type: c.Type,
			}
			break
		}
	}
	if cp == nil {
		http.Error(w, "Checkpoint no encontrado", http.StatusNotFound)
		return
	}

	// Validar distancia
	dist := haversine(input.Lat, input.Lng, cp.Lat, cp.Lng)
	if dist > checkpointRadius {
		http.Error(w, "Fuera de rango del checkpoint", http.StatusBadRequest)
		return
	}

	// Guardar checkin
	if err := repository.CreateCheckin(claims.UserID, eventID, checkpointID, input.Lat, input.Lng); err != nil {
		http.Error(w, "Error registrando checkin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"checkpoint": cp.Name,
		"status":     "ok",
		"distance_m": dist,
		"message":    "Checkpoint validado correctamente",
	})
}
