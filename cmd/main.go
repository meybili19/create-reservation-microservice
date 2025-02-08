package main

import (
	"log"
	"net/http"

	"github.com/meybili19/create-reservation-microservice/config"
	"github.com/meybili19/create-reservation-microservice/routes"
)

func main() {
	databases, err := config.InitDatabases()
	if err != nil {
		log.Fatalf("Error initializing databases: %v", err)
	}
	defer func() {
		for _, db := range databases {
			db.Close()
		}
	}()
	log.Println("All databases connected successfully!")

	http.HandleFunc("/reservations", routes.CreateReservationHandler(databases))
	log.Println("Server running on port 4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
