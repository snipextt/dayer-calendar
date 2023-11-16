package storage

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var kafkaWriter *kafka.Writer

func connectToKafka() {
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	mechanism, err := scram.Mechanism(scram.SHA256, username, password)
	if err != nil {
		log.Fatalln(err)
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"fancy-adder-10821-eu1-kafka.upstash.io:9092"},
		Topic:   "timedoctorReport",
		Dialer:  dialer,
	})

	kafkaWriter = w
}

func KafkaWriter() *kafka.Writer {
	return kafkaWriter
}
