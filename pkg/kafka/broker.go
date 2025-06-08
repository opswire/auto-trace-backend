package kafka

import (
	"github.com/segmentio/kafka-go"
	"os"
)

func NewBroker() kafka.Writer {
	broker := os.Getenv("KAFKA_BROKER")
	writer := kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "payments",
		Balancer: &kafka.LeastBytes{},
	}

	return writer
}
