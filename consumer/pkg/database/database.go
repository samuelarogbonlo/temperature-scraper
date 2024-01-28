package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

type CityTemperature struct {
	City        string
	Temperature string
	Time        string
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
		return err
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

func (db *DB) GetCityTemperature(city string) ([]CityTemperature, error) {
	var temperatures []CityTemperature

	query := `SELECT city, temperature, time FROM city_temperatures WHERE city = $1`
	rows, err := db.Query(query, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp CityTemperature
		if err := rows.Scan(&temp.City, &temp.Temperature, &temp.Time); err != nil {
			return nil, err
		}
		temperatures = append(temperatures, temp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return temperatures, nil
}
