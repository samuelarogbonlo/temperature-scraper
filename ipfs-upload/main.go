package main

import (
	"log"
	"os"
	"fmt"
	"encoding/json"
	"time"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/sqltocsv"
	"github.com/ipfs/go-ipfs-api"
)

type CityTemperature struct {
	City        string
	Temperature string
	Time        string
}

type WeatherAggregate struct {
    City            string
    Date            string
    AvgTemperature  float64
    MedianTemperature float64
    MinTemperature  float64
    MaxTemperature  float64
}

func uploadToIPFS(filename string)(string, error){
	sh := shell.NewShell("localhost:5001")
	cid, err := sh.AddDir(filename)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func main(){
	if len(os.Args) < 3 {
        log.Fatal("Usage: go run script.go [city] [date in YYYY-MM-DD]")
    }

    city := os.Args[1]
    date := os.Args[2]

    // Validate date format
    _, date_err := time.Parse("2006-01-02", date)
    if date_err != nil {
        log.Fatal("Invalid date format. Please use YYYY-MM-DD.")
    }
	env_err := godotenv.Load()
	if env_err != nil{
		log.Fatalf("Error loading .env file: %s", env_err)
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)
	// dsn := os.ExpandEnv("host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=5432 sslmode=${SSL_MODE}")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	var result []map[string]interface{}

	tx := db.Find(&CityTemperature{}).Scan(&result)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes))

	//var results []WeatherAggregate

	query := `
		SELECT city,
			date_trunc('day', time) AS date,
			AVG(NULLIF(REGEXP_REPLACE(temperature, '[^0-9.-]', '', 'g'), '')::float) AS avg_temperature,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY NULLIF(REGEXP_REPLACE(temperature, '[^0-9.-]', '', 'g'), '')::float) AS median_temperature,
			MIN(NULLIF(REGEXP_REPLACE(temperature, '[^0-9.-]', '', 'g'), '')::float) AS min_temperature,
			MAX(NULLIF(REGEXP_REPLACE(temperature, '[^0-9.-]', '', 'g'), '')::float) AS max_temperature
		FROM 
			city_temperatures
		WHERE
			city = ? AND date_trunc('day', time) = ?
		GROUP BY 
			city, date_trunc('day', time)
		ORDER BY city, date;
		`

	// Execute the query
    rows, err := db.Raw(query, city, date).Rows()
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

	if !rows.Next() {
        log.Fatalf("No data found for city '%s' on date '%s'", city, date)
    }

    // Using sqltocsv to write directly to a file
	filename := fmt.Sprintf("temp_data_points_%s_%s.csv", city, date)
    err = sqltocsv.WriteFile(filename, rows)
    if err != nil {
        log.Fatal("Unable to write to file:", err)
    }
    log.Println("Export completed successfully")
	cid, err := uploadToIPFS(filename)
	if err != nil {
        log.Fatal("Failed to upload to IPFS: ", err)
    }

    fmt.Printf("File uploaded to IPFS with CID: %s\n", cid)
    fmt.Printf("Publicly accessible link (via IPFS gateway): https://ipfs.io/ipfs/%s\n", cid)
}
