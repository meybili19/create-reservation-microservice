package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB establece una conexión a una base de datos MySQL.
func ConnectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying connection to database: %w", err)
	}
	return db, nil
}

// InitDatabases inicializa todas las conexiones a las bases de datos requeridas.
func InitDatabases() (map[string]*sql.DB, error) {
	// DSNs para cada base de datos
	databases := map[string]string{
		"users":        "root:Mtoi2002.@tcp(localhost:3306)/UsersDB",
		"cars":         "root:Mtoi2002.@tcp(localhost:3306)/CarDB",
		"parkinglots":  "root:Mtoi2002.@tcp(localhost:3306)/ParkingLotDB",
		"reservations": "root:Mtoi2002.@tcp(localhost:3306)/ReservationDB",
	}

	// Inicializamos las conexiones
	connections := make(map[string]*sql.DB)
	for name, dsn := range databases {
		db, err := ConnectDB(dsn)
		if err != nil {
			return nil, fmt.Errorf("error connecting to %s: %w", name, err)
		}
		connections[name] = db
	}
	return connections, nil
}
