package kafkaProducer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

type CityTemperatureData struct {
	City        string    `json:"city"`
	Temperature string    `json:"temperature"`
	Time        time.Time `json:"time"`
}

var (
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

const (
	KafkaTopic          = "notifications"
	ProducerPort        = ":4000"
)

func SetupProducer() (sarama.SyncProducer, error) {
	kafkaAddress := os.Getenv("KAFKA_SERVER_ADDRESS") 
	if kafkaAddress == "" {
		return nil, fmt.Errorf("kafka address not set in environment variable KAFKA_SERVER_ADDRESS")
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaAddress}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return producer, nil
}


func SendKafkaMessage(producer sarama.SyncProducer, cityTempData CityTemperatureData) error {
	dataJSON, err := json.Marshal(cityTempData)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(cityTempData.City), // City as the key
		Value: sarama.StringEncoder(dataJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}
