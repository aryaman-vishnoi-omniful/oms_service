package services

import (
	"context"
	"encoding/json"
	"fmt"
	kafka_producer "oms_service/kafka"
	"oms_service/orders/requests"
	"os"
	"strconv"
	"time"

	// "oms_service/orders/listners"
	// "time"

	// "fmt"
	"log"

	// "github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/pubsub"
	"github.com/omniful/go_commons/sqs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var NewProducer = &sqs.Publisher{}

// var Queue_instance =&sqs.Queue{}

func SetProducer(ctx context.Context, queue *sqs.Queue) {
	NewProducer = sqs.NewPublisher(queue)

}
func SendMessage(ctx context.Context, message *sqs.Message) {

	err := NewProducer.Publish(ctx, message)
	if err != nil {
		log.Fatal("did not publish", err)
	}

}

//	type CSVUploadRequest struct {
//		FilePath string
//		UserID   string
//	}
func PushCreateOrderMessageToKafka(ctx context.Context, order *requests.Order) error {
	// Marshal the order's items to JSON. (Change to marshal the entire order if desired.)
	msg, err := json.Marshal(order.OrderItems)
	if err != nil {
		return fmt.Errorf("unable to marshal order items: %w", err)
	}

	headers := map[string]string{
		"event": "order_created",
	}

	// Use OrderNo as the Kafka partition key. (Alternatively, you can use order.ID.String())
	key := order.OrderNo

	p := kafka_producer.Get()
	publishErr := p.Publish(ctx, &pubsub.Message{
		Topic:   config.GetString(ctx, "consumers.orders.topic"),
		Value:   msg,
		Key:     key,
		Headers: headers,
	})
	if publishErr != nil {
		return fmt.Errorf("could not publish order %s to Kafka: %w", order.OrderNo, publishErr)
	}

	return nil
}
func ExtractFromCsv(filePath string) ([]*requests.Order, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Map to group items by order_no and customer_name
	ordermap := make(map[string]*requests.Order)

	// Initialize the CSV reader (using your CSV package and options)
	Csv, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),
		csv.WithSource(csv.Local),
		csv.WithLocalFileInfo(filePath),
		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}
	if err = Csv.InitializeReader(context.TODO()); err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}

	// Process the records and group them by order_no and customer_name
	for !Csv.IsEOF() {
		records, err := Csv.ReadNextBatch()
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV batch: %w", err)
		}

		fmt.Println("Processing records:")
		fmt.Println(records)
		for _, record := range records {
			orderNo := record[0]      // order_no
			customerName := record[1] // customer_name
			skuID := record[2]        // sku_id
			quantityStr := record[3]  // quantity

			// Convert quantity to integer
			quantity, err := strconv.Atoi(quantityStr)
			if err != nil {
				return nil, fmt.Errorf("invalid quantity %s: %v", quantityStr, err)
			}

			// Group by a key (order_no-customer_name)
			orderKey := fmt.Sprintf("%s-%s", orderNo, customerName)
			order, exists := ordermap[orderKey]
			if !exists {
				// Create a new order if it doesn't exist
				now := primitive.NewDateTimeFromTime(time.Now())
				order = &requests.Order{
					ID:           primitive.NewObjectID(),
					CustomerName: customerName,
					OrderNo:      orderNo,
					OrderItems:   []requests.OrderItem{},
					Status:       "on_hold",
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				ordermap[orderKey] = order
			}

			// Create an OrderItem for this record
			orderItem := requests.OrderItem{
				OrderID: orderNo, // or order.ID.String(), depending on your schema
				SKUID:   skuID,
				Quantity: quantity,
			}
			order.OrderItems = append(order.OrderItems, orderItem)
		}
	}

	// Convert the map of orders into a slice
	var orders []*requests.Order
	for _, order := range ordermap {
		orders = append(orders, order)
	}

	fmt.Println("Final orders:")
	for _, order := range orders {
		fmt.Printf("Order No: %s, Customer: %s, Total Items: %d\n", order.OrderNo, order.CustomerName, len(order.OrderItems))
	}

	return orders, nil
}
func ParseCSV(filePath string, ctx context.Context) {
	orders, err := ExtractFromCsv(filePath)
	if err != nil {
		fmt.Printf("\nFailed to parse CSV (%s): %v\n", filePath, err)
		return
	}
	
	// Publish each order to Kafka
	for _, order := range orders {
		if err := PushCreateOrderMessageToKafka(ctx, order); err != nil {
			fmt.Printf("Failed to publish order %s: %v\n", order.OrderNo, err)
		} else {
			fmt.Printf("Published order %s successfully\n", order.OrderNo)
		}
	}
}

func ConvertControllerReqToServiceReqParseCsv(ctx context.Context, CreateBulkCsv *requests.CSVUploadRequest) (string, error) {

	message := &sqs.Message{
		GroupId:         "csv-163",
		Value:           []byte(CreateBulkCsv.FilePath),
		ReceiptHandle:   "order-csv3",
		DeduplicationId: "dedup-490",
		DelayDuration:   0,
		Headers: map[string]string{
			"customer_id": CreateBulkCsv.CustomerId,
		},
	}
	SendMessage(ctx, message)

	log.Println("Message sent successfully to SQS")
	return "", nil

}
