package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"

	"dia/interview/temp-scraper/pkg/kafkaProducer"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func getTemperature(city string) (string, error) {
	route := fmt.Sprintf("http://wttr.in/%s?format=%%t", city)

	timer := prometheus.NewTimer(httpDuration.WithLabelValues(route))
	defer timer.ObserveDuration()

	resp, err := http.Get(route)
	if err != nil {
		totalRequests.WithLabelValues(route).Inc()
		responseStatus.WithLabelValues("error").Inc()
		return "", err
	}
	defer resp.Body.Close()

	totalRequests.WithLabelValues(route).Inc()
	responseStatus.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {

	producer, err := kafkaProducer.SetupProducer()
	if err != nil {
		kafkaProducer.ErrorLog.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	kafkaProducer.InfoLog.Println("Kafka producer initialized successfully")
	defer producer.Close()


	cronJob := cron.New()
	kafkaProducer.InfoLog.Println("Starting Temperature Scheduler")
	kafkaProducer.InfoLog.Printf("Starting HTTP server on port %s", kafkaProducer.ProducerPort)
	cronJob.AddFunc("@hourly", func() {
		cities := []string{"Zurich", "London", "Miami", "Tokyo", "Singapore"}
		for _, city := range cities {
			temp, err := getTemperature(city)
			if err != nil {
				kafkaProducer.ErrorLog.Printf("Could not get temperature for %s: %s\n", city, err)
				continue
			}

			cityTempData := kafkaProducer.CityTemperatureData{
				City:        city,
				Temperature: temp,
				Time:        time.Now(),
			}

			err = kafkaProducer.SendKafkaMessage(producer, cityTempData)
			if err != nil {
				kafkaProducer.ErrorLog.Printf("Failed to send Kafka message for %s: %s\n", city, err)
				continue
			}
			kafkaProducer.InfoLog.Printf("Kafka message sent for %s", city)
		}
	})
	cronJob.Start()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "alive")
	})
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		kafkaProducer.ErrorLog.Fatal(http.ListenAndServe(":8080", nil))
	}()

	select {}
}
