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
	// 1️⃣ VALIDAR QUE EL VEHÍCULO EXISTE
	carID := int(reservation["car_id"].(float64))
	vehicleServiceURL := fmt.Sprintf("%s/%d", os.Getenv("VEHICLE_SERVICE_URL"), carID)

	resp, err := http.Get(vehicleServiceURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("vehicle not found")
	}
	defer resp.Body.Close()

	// Leer respuesta JSON del microservicio de vehículos
	var vehicleData map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &vehicleData)

	// 2️⃣ OBTENER EL `user_id` ASOCIADO AL VEHÍCULO
	var userID int
	if uid, ok := vehicleData["user_id"].(float64); ok {
		userID = int(uid)
	} else if uid, ok := vehicleData["userId"].(float64); ok {
		userID = int(uid)
	} else {
		return errors.New("invalid user_id received from vehicle service")
	}
	reservation["user_id"] = userID

	// 3️⃣ VALIDAR QUE EL PARQUEADERO EXISTE Y TIENE CAPACIDAD
	parkingLotID := int(reservation["parking_lot_id"].(float64))

	// Validar que el parqueadero existe y tiene capacidad suficiente
	if err := parkinglot.CheckParkingLotAvailability(parkingLotID); err != nil {
		return err
	}

	// 4️⃣ CALCULAR EL PRECIO DE LA RESERVA
	totalAmount, err := utils.CalculatePrice(reservation)
	if err != nil {
		return err
	}
	reservation["total_amount"] = totalAmount

	// 5️⃣ INSERTAR LA RESERVA EN LA BASE DE DATOS
	reservation["status"] = "Pending"

	err = repositories.CreateReservation(databases["reservations"], reservation)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %v", err)
	}

	// 6️⃣ DISMINUIR CAPACIDAD DEL PARQUEADERO (AHORA USANDO PUT)
	if err := parkinglot.DecreaseParkingLotCapacity(parkingLotID); err != nil {
		return fmt.Errorf("reservation created but failed to update parking lot capacity: %w", err)
	}

	return nil
}
