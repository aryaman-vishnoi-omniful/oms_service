package listners

import (
	"context"
	"encoding/json"
	"oms_service/orders/requests"

	// "oms-service/r"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

// Implement message handler
type MessageHandler struct{}

// Process implements pubsub.IPubSubMessageHandler.
func (h *MessageHandler) Process(ctx context.Context, message *pubsub.Message) error {
	log.Printf("Received message: %s", string(message.Value))

	var orders []requests.KafkaResponseOrderMessage 
	err := json.Unmarshal(message.Value, &orders)
	if err != nil {
		log.Printf("Failed to parse Kafka message: %v", err)
		return err
	}

	log.Printf("Parsed Kafka Order Messages: %+v", orders)

	// Process each order
	for _, order := range orders {
		log.Printf("Processing Order: %+v", order)
		
	}

	return nil
}


func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
	// Process message
	return nil
}

// Initialize Kafka Consumer
func InitializeKafkaConsumer(ctx context.Context) {
	consumer := kafka.NewConsumer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithConsumerGroup("my-consumer-group"),
		kafka.WithClientID("my-consumer"),
		kafka.WithKafkaVersion("2.8.1"),
		kafka.WithRetryInterval(time.Second),
	)

	handler := &MessageHandler{}
	consumer.RegisterHandler(config.GetString(ctx, "consumers.orders.topic"), handler)
	consumer.Subscribe(ctx)
}
