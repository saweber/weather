package main

import (
	"fmt"

	"github.com/IBM/sarama"
)

// Sarama configuration options
var (
	brokers = "localhost:9092"
	version = sarama.DefaultVersion.String()
	group   = "etl"
	topics  = "raw-weather-reports"
)

func main() {
	// create repository to save to postgres
	r := &Repository{}
	r.Init()

	// create kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Printf("error creating Kafka producer: %v", err)
		return
	}
	defer producer.Close()

	e := ETL{
		r:        r,
		producer: producer,
	}
	e.ReadFromKafkaAndProcess()
}
