package models

type Reservation struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	CarID        int     `json:"car_id"`
	ParkingLotID int     `json:"parking_lot_id"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
}
