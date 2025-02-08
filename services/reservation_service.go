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
	parkingLotURL := fmt.Sprintf("%s/%d", os.Getenv("PARKINGLOT_SERVICE_URL"), parkingLotID)
	resp, err = http.Get(parkingLotURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("parking lot not found")
	}
	defer resp.Body.Close()

	capacityURL := fmt.Sprintf("%s/%d", os.Getenv("PARKINGLOT_SERVICE_CAPACITY_URL"), parkingLotID)
	resp, err = http.Get(capacityURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("error checking parking lot capacity")
	}
	defer resp.Body.Close()

	var capacityData map[string]interface{}
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &capacityData)

	if capacity, ok := capacityData["capacity"].(float64); ok && int(capacity) == 0 {
		return errors.New("parking lot has no available capacity")
	}

	totalAmount, err := utils.CalculatePrice(reservation)
	if err != nil {
		return err
	}
	reservation["total_amount"] = totalAmount

	// Set default status
	reservation["status"] = "Pending"

	// Create the reservation
	err = repositories.CreateReservation(databases["reservations"], reservation)
	if err != nil {
		return err
	}

	// Decrease parking lot capacity after successful reservation
	updateCapacityURL := fmt.Sprintf("%s/%d/decrease", os.Getenv("PARKINGLOT_SERVICE_URL"), parkingLotID)
	req, err := http.NewRequest("PUT", updateCapacityURL, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("failed to update parking lot capacity")
	}
	defer resp.Body.Close()

	return nil
}
