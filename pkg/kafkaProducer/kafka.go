package kafka

import (
    "encoding/json"
    "github.com/IBM/sarama"
    "time"
    "log"
    "os"
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
    KafkaEnvVariable = "KAFKA_SERVER_ADDRESS" // Name of the environment variable
    DefaultKafkaAddress = "localhost:9092"    // Default address if env variable is not set
    KafkaTopic = "notifications"
    ProducerPort = ":8080" 
)

func SetupProducer() (sarama.SyncProducer, error) {
    kafkaAddress := os.Getenv(KafkaEnvVariable)
    if kafkaAddress == "" {
        kafkaAddress = DefaultKafkaAddress
    }

    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    producer, err := sarama.NewSyncProducer([]string{kafkaAddress}, config)
    if err != nil {
        return nil,  err
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
