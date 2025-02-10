package utils

import (
	"fmt"
	"math"
	"time"
)

const pricePerHour = 0.50
const pricePerDay = 5.00
const monthlyPrice = 60.00
const maxDays = 30

// CalculatePrice calcula el costo total de la reserva
func CalculatePrice(reservation map[string]interface{}) (float64, error) {
	startDateStr, ok := reservation["start_date"].(string)
	if !ok {
		return 0, fmt.Errorf("start_date must be a valid string")
	}

	endDateStr, ok := reservation["end_date"].(string)
	if !ok {
		return 0, fmt.Errorf("end_date must be a valid string")
	}

	layout := "2006-01-02 15:04:05"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing start_date: %v", err)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing end_date: %v", err)
	}

	duration := endDate.Sub(startDate)
	totalHours := duration.Hours()
	totalDays := int(totalHours / 24)

	if totalHours > float64(maxDays*24) {
		return 0, fmt.Errorf("cannot reserve for more than 30 days")
	}

	var totalAmount float64
	if totalHours == float64(maxDays*24) {
		totalAmount = monthlyPrice
	} else if totalDays > 0 {
		totalAmount = float64(totalDays) * pricePerDay
		remainingHours := totalHours - float64(totalDays)*24
		if remainingHours > 0 {
			totalAmount += remainingHours * pricePerHour
		}
	} else {
		if totalHours < 1 {
			if duration.Minutes() >= 1 {
				totalAmount = duration.Minutes() * (pricePerHour / 60)
			} else {
				totalAmount = duration.Seconds() * (pricePerHour / 3600)
			}
		} else {
			totalAmount = totalHours * pricePerHour
		}
	}

	return math.Round(totalAmount*100) / 100, nil
}
