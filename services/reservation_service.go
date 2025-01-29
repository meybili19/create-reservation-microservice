package services

import (
	"database/sql"
	"fmt"

	"github.com/meybili19/create-reservation-microservice/repositories"
)

func CreateReservationService(databases map[string]*sql.DB, reservation map[string]interface{}) error {
	if databases["users"] == nil || databases["cars"] == nil || databases["parkinglots"] == nil || databases["reservations"] == nil {
		return fmt.Errorf("una o más bases de datos no están inicializadas correctamente")
	}

	userID, ok := reservation["user_id"].(float64)
	if !ok {
		return fmt.Errorf("user_id debe ser un número válido")
	}

	carID, ok := reservation["car_id"].(float64)
	if !ok {
		return fmt.Errorf("car_id debe ser un número válido")
	}

	parkingLotID, ok := reservation["parking_lot_id"].(float64)
	if !ok {
		return fmt.Errorf("parking_lot_id debe ser un número válido")
	}

	if err := repositories.CheckForeignKey(databases["users"], "Users", int(userID)); err != nil {
		return err
	}

	if err := repositories.CheckForeignKey(databases["cars"], "Cars", int(carID)); err != nil {
		return err
	}

	if err := repositories.CheckForeignKey(databases["parkinglots"], "ParkingLot", int(parkingLotID)); err != nil {
		return err
	}

	if _, exists := reservation["status"]; !exists {
		reservation["status"] = "Pending"
	}

	return repositories.CreateReservation(databases["reservations"], reservation)
}
