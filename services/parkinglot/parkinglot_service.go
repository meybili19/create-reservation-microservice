package parkinglot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Verifica que el parqueadero existe y tiene espacio disponible
func CheckParkingLotAvailability(parkingLotID int) error {
	parkingLotURL := fmt.Sprintf("%s/%d", os.Getenv("PARKINGLOT_SERVICE_URL"), parkingLotID)
	resp, err := http.Get(parkingLotURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("parking lot not found")
	}
	defer resp.Body.Close()

	// Consultar la capacidad del parqueadero
	capacityURL := fmt.Sprintf("%s/%d", os.Getenv("PARKINGLOT_SERVICE_CAPACITY_URL"), parkingLotID)
	resp, err = http.Get(capacityURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("error checking parking lot capacity")
	}
	defer resp.Body.Close()

	var capacityData map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &capacityData)

	// Validar que haya espacio disponible
	if capacity, ok := capacityData["capacity"].(float64); ok && int(capacity) == 0 {
		return errors.New("parking lot has no available capacity")
	}

	return nil
}

func DecreaseParkingLotCapacity(parkingLotID int) error {
	// Construir la URL correcta con el ID del parqueadero
	decreaseCapacityURL := fmt.Sprintf("%s/%d", os.Getenv("PARKINGLOT_SERVICE_DISMINUYE_URL"), parkingLotID)

	req, err := http.NewRequest("PUT", decreaseCapacityURL, nil) // 🔹 Cambio de POST a PUT
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to %s: %v", decreaseCapacityURL, err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update parking lot capacity. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	return nil
}
