package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/meybili19/create-reservation-microservice/services"
)

func CreateReservationHandler(databases map[string]*sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reservation map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := services.CreateReservationService(databases, reservation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Reservation created successfully",
		})
	}
}
