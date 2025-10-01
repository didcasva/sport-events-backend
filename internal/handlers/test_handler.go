package handlers

import (
	"encoding/json"
	"net/http"

	"sport-events-backend/internal/config"
)

func TestDBHandler(w http.ResponseWriter, r *http.Request) {
	type Row struct {
		ID   int    `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}

	var rows []Row
	err := config.DB.Select(&rows, "SELECT * FROM test_connection")
	if err != nil {
		http.Error(w, "Error al consultar BD: "+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(rows)
}
