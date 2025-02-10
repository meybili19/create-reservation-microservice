package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/meybili19/create-reservation-microservice/repositories"
	"github.com/meybili19/create-reservation-microservice/services/parkinglot"
	"github.com/meybili19/create-reservation-microservice/utils"
)

func CreateReservationService(databases map[string]*sql.DB, reservation map[string]interface{}) error {

	carID := int(reservation["car_id"].(float64))
	vehicleServiceURL := fmt.Sprintf("%s/%d", os.Getenv("VEHICLE_SERVICE_URL"), carID)

	resp, err := http.Get(vehicleServiceURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("vehicle not found")
	}
	defer resp.Body.Close()

	var vehicleData map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &vehicleData)

	var userID int
	if uid, ok := vehicleData["user_id"].(float64); ok {
		userID = int(uid)
	} else if uid, ok := vehicleData["userId"].(float64); ok {
		userID = int(uid)
	} else {
		return errors.New("invalid user_id received from vehicle service")
	}
	reservation["user_id"] = userID

	parkingLotID := int(reservation["parking_lot_id"].(float64))

	if err := parkinglot.CheckParkingLotAvailability(parkingLotID); err != nil {
		return err
	}

	totalAmount, err := utils.CalculatePrice(reservation)
	if err != nil {
		return err
	}
	reservation["total_amount"] = totalAmount

	reservation["status"] = "Pending"

	err = repositories.CreateReservation(databases["reservations"], reservation)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %v", err)
	}

	if err := parkinglot.DecreaseParkingLotCapacity(parkingLotID); err != nil {
		return fmt.Errorf("reservation created but failed to update parking lot capacity: %w", err)
	}

	return nil
}
