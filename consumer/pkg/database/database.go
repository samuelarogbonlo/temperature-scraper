package database

import (
    "database/sql"
    "fmt"
    _ "log"
	
    _ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
}

type CityTemperature struct {
    City        string
    Temperature string
    Time        string // use string for simplicity
}


func (db *DB) InitializeSchema() error {
    query := `
    CREATE TABLE IF NOT EXISTS city_temperatures (
        id SERIAL PRIMARY KEY,
        city VARCHAR(255),
        temperature VARCHAR(50),
        time TIMESTAMP
    );
    `

    if _, err := db.Exec(query); err != nil {
        return fmt.Errorf("error creating city_temperatures table: %w", err)
    }
    return nil
}


func New(connectionString string) (*DB, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return &DB{db}, nil
}

func (db *DB) SaveCityTemperature(data CityTemperature) error {
    query := `INSERT INTO city_temperatures (city, temperature, time) VALUES ($1, $2, $3)`
    _, err := db.Exec(query, data.City, data.Temperature, data.Time)
    if err != nil {
        return err
    }
    return nil
}
