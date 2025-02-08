package listners

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"oms_service/orders/requests"

	// "oms-service/r"
	"time"

	"github.com/omniful/go_commons/config"
	interservice_client "github.com/omniful/go_commons/interservice-client"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"

	"github.com/omniful/go_commons/http"
)

// Implement message handler
type MessageHandler struct{}

func ValidateInventory(ctx context.Context, order requests.KafkaResponseOrderMessage) error {

	log.Printf("Validating inventory for order ID: %s \n", order.OrderID)

	// client := &http2.Client{}
	// url := "http://localhost:8081/api/v1/orders/validate_inventory"

	config := interservice_client.Config{
		ServiceName: "order-service",
		BaseURL:     "http://localhost:8081/wms/v1/",
		Timeout:     5 * time.Second,
	}
	client, err := interservice_client.NewClientWithConfig(config)
	if err != nil {
		return err
	}

	url := config.BaseURL + "inventory/deduct"
	bodyBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req := &http.Request{
		Url:     url, // Use configured URL
		Body:    bytes.NewReader(bodyBytes),
		Timeout: 7 * time.Second,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	_, intersvcErr := client.Put(req, "/")
	if intersvcErr.StatusCode.Is4xx() {
		fmt.Println("inventory validation failed after  interservice call to wms-service ")
	} else {
		fmt.Println("Inventory validation successful")
	}

	return nil
}

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
		err = ValidateInventory(ctx, order)
		if err != nil {
		log.WithError(err).Error("Inventory validation failed \n")
		return err
	}
	// return nil
		
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
