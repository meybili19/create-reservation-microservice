package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

func CheckForeignKey(db *sql.DB, table string, id int) error {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = ?", table)
	row := db.QueryRow(query, id)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s with ID %d not found", table, id)
		}
		return err
	}
	return nil
}

func CreateReservation(db *sql.DB, reservation map[string]interface{}) error {
	query := `INSERT INTO Reservations 
		(user_id, car_id, parking_lot_id, start_date, end_date, status, total_amount) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, reservation["user_id"], reservation["car_id"], reservation["parking_lot_id"],
		reservation["start_date"], reservation["end_date"], reservation["status"], reservation["total_amount"])
	return err
}
