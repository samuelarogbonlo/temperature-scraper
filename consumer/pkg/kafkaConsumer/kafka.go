package kafkaConsumer

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
    "os"
	"github.com/IBM/sarama"
    "dia/interview/consumer/pkg/database"
)

type CityTemperatureData struct {
    City        string    `json:"city"`
    Temperature string    `json:"temperature"`
    Time        time.Time `json:"time"`
}

const (
    KafkaServerAddress = "localhost:9092"
    KafkaTopic         = "notifications"
    ConsumerGroup      = "notifications-group"
)

var (
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

// DataStore stores the received messages
type DataStore struct {
    mu   sync.Mutex
    data []CityTemperatureData
}

func NewDataStore() *DataStore {
    return &DataStore{
        data: []CityTemperatureData{},
    }
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
    ready chan bool
    Store *DataStore
    DB    *database.DB
}


func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
    close(consumer.ready)
    return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var data CityTemperatureData
        if err := json.Unmarshal(msg.Value, &data); err != nil {
            ErrorLog.Printf("Error unmarshalling data: %v", err)
            continue
        }

        // Convert data to the database struct and save
        dbData := database.CityTemperature{
            City:        data.City,
            Temperature: data.Temperature,
            Time:        data.Time.Format(time.RFC3339),
        }
        if err := consumer.DB.SaveCityTemperature(dbData); err != nil {
            ErrorLog.Printf("Failed to save data to database: %v", err)
            continue
        }
        InfoLog.Printf("Temperature Data for %s saved successfully", data.City)
        sess.MarkMessage(msg, "")
    }
    return nil
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
    config := sarama.NewConfig()
    config.Version = sarama.MaxVersion

    // Start consuming from the oldest message
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    consumerGroup, err := sarama.NewConsumerGroup([]string{KafkaServerAddress}, ConsumerGroup, config)
    if err != nil {
        return nil, err
    }

    return consumerGroup, nil
}

func NewConsumer(store *DataStore, db *database.DB) *Consumer {
    return &Consumer{
        ready: make(chan bool),
        Store: store,
        DB:    db,
    }
}

func SetupConsumerGroup(ctx context.Context, consumer *Consumer) {
    consumerGroup, err := initializeConsumerGroup()
    if err != nil {
        ErrorLog.Fatalf("Failed to initialize consumer group: %v", err)
    }
    defer consumerGroup.Close()

    for {
        if err := consumerGroup.Consume(ctx, []string{KafkaTopic}, consumer); err != nil {
            ErrorLog.Printf("Error from consumer: %v", err)
            time.Sleep(5 * time.Second) // retry after a delay
        }
        select {
        case <-ctx.Done():
            ErrorLog.Println("Consumer context cancelled, exiting consumer loop")
            return
        default:
            <-consumer.ready
        }
    }
}
