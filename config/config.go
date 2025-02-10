package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

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

	// 🔹 Construir manualmente los DSN con las variables de entorno
	databases := map[string]string{
		"parkinglots": fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			os.Getenv("DB_PARKINGLOTS_USER"),
			os.Getenv("DB_PARKINGLOTS_PASSWORD"),
			os.Getenv("DB_PARKINGLOTS_HOST"),
			os.Getenv("DB_PARKINGLOTS_NAME"),
		),
		"reservations": fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			os.Getenv("DB_RESERVATIONS_USER"),
			os.Getenv("DB_RESERVATIONS_PASSWORD"),
			os.Getenv("DB_RESERVATIONS_HOST"),
			os.Getenv("DB_RESERVATIONS_NAME"),
		),
	}

	connections := make(map[string]*sql.DB)
	for name, dsn := range databases {
		if dsn == "" {
			return nil, fmt.Errorf("missing DSN for %s", name)
		}
		db, err := ConnectDB(dsn)
		if err != nil {
			return nil, fmt.Errorf("error connecting to %s: %w", name, err)
		}
		connections[name] = db
	}
	return connections, nil
}
