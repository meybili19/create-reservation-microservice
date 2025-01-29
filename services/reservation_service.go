package services

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/meybili19/create-reservation-microservice/repositories"
)

const pricePerHour = 0.50
const pricePerDay = 5.00   // Price per day
const monthlyPrice = 60.00 // Price for 30 days
const maxDays = 30         // Maximum number of days for a reservation

func CreateReservationService(databases map[string]*sql.DB, reservation map[string]interface{}) error {
	if databases["users"] == nil || databases["cars"] == nil || databases["parkinglots"] == nil || databases["reservations"] == nil {
		return fmt.Errorf("one or more databases are not initialized correctly")
	}

	userID, ok := reservation["user_id"].(float64)
	if !ok {
		return fmt.Errorf("user_id must be a valid number")
	}

	carID, ok := reservation["car_id"].(float64)
	if !ok {
		return fmt.Errorf("car_id must be a valid number")
	}

	parkingLotID, ok := reservation["parking_lot_id"].(float64)
	if !ok {
		return fmt.Errorf("parking_lot_id must be a valid number")
	}

	startDateStr, ok := reservation["start_date"].(string)
	if !ok {
		return fmt.Errorf("start_date must be a valid string")
	}

	endDateStr, ok := reservation["end_date"].(string)
	if !ok {
		return fmt.Errorf("end_date must be a valid string")
	}

	// Convert dates to time.Time type
	layout := "2006-01-02 15:04:05" // Standard date and time format for MySQL
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return fmt.Errorf("error parsing start_date: %v", err)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return fmt.Errorf("error parsing end_date: %v", err)
	}

	// Calculate total duration in hours and days
	duration := endDate.Sub(startDate)
	totalHours := duration.Hours()
	totalDays := int(totalHours / 24)

	// Validate if the reservation exceeds 30 days
	if totalHours > float64(maxDays*24) {
		return fmt.Errorf("cannot reserve for more than 30 days. Please renew your reservation after one month")
	}

	// Logic to calculate the total amount
	var totalAmount float64

	// If the reservation is exactly 30 days, charge $60
	if totalHours == float64(maxDays*24) {
		totalAmount = monthlyPrice
	} else if totalDays > 0 {
		// If the reservation is for 1 or more days, charge $5 per day
		totalAmount = float64(totalDays) * pricePerDay

		// Calculate additional hours after the full days
		remainingHours := totalHours - float64(totalDays)*24
		if remainingHours > 0 {
			totalAmount += remainingHours * pricePerHour
		}
	} else {
		// If the reservation is less than 1 day, calculate by hours, minutes, or seconds
		if totalHours < 1 {
			if duration.Minutes() >= 1 {
				totalAmount = duration.Minutes() * (pricePerHour / 60) // Charge by minute
			} else {
				totalAmount = duration.Seconds() * (pricePerHour / 3600) // Charge by second
			}
		} else {
			// Charge by hours
			totalAmount = totalHours * pricePerHour
		}
	}

	// Round the total amount to 2 decimal places
	totalAmount = math.Round(totalAmount*100) / 100

	// Assign the total amount to the reservation
	reservation["total_amount"] = totalAmount

	// Check foreign keys
	if err := repositories.CheckForeignKey(databases["users"], "Users", int(userID)); err != nil {
		return err
	}

	if err := repositories.CheckForeignKey(databases["cars"], "Cars", int(carID)); err != nil {
		return err
	}

	if err := repositories.CheckForeignKey(databases["parkinglots"], "ParkingLot", int(parkingLotID)); err != nil {
		return err
	}

	// Set the reservation status if not defined
	if _, exists := reservation["status"]; !exists {
		reservation["status"] = "Pending"
	}

	// Create the reservation in the database
	return repositories.CreateReservation(databases["reservations"], reservation)
}
