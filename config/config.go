package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

func InitDatabases() (map[string]*sql.DB, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	databases := map[string]string{
		"users":        os.Getenv("DB_USERS_DSN"),
		"cars":         os.Getenv("DB_CARS_DSN"),
		"parkinglots":  os.Getenv("DB_PARKINGLOTS_DSN"),
		"reservations": os.Getenv("DB_RESERVATIONS_DSN"),
	}

	for name, dsn := range databases {
		if dsn == "" {
			return nil, fmt.Errorf("missing DSN for %s", name)
		}
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
