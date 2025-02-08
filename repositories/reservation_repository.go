package repositories

import (
	"database/sql"
)

func CreateReservation(db *sql.DB, reservation map[string]interface{}) error {
	query := `INSERT INTO Reservations (user_id, car_id, parking_lot_id, start_date, end_date, status, total_amount) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, reservation["user_id"], reservation["car_id"], reservation["parking_lot_id"], reservation["start_date"], reservation["end_date"], reservation["status"], reservation["total_amount"])
	return err
}
