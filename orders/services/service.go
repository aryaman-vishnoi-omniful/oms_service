package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	// "net/http"
	kafka_producer "oms_service/kafka"
	"oms_service/orders/requests"

	// "oms_service/orders/services"
	"oms_service/repository"
	"os"
	"strconv"
	"time"

	"github.com/omniful/go_commons/http"

	// "oms_service/orders/listners"
	// "time"

	// "fmt"
	"log"

	// "github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/csv"
	interservice_client "github.com/omniful/go_commons/interservice-client"
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
	fmt.Println("pushed to queue")

}

//	type CSVUploadRequest struct {
//		FilePath string
//		UserID   string
//	}

// const BaseURL = "http://localhost:8081/wms/v1"

func validateSKU(skuID int) bool {
	skustr := strconv.Itoa(int(skuID))

	config := interservice_client.Config{
		ServiceName: "order-service",
		BaseURL:     "http://localhost:8081/wms/v1/getSkuById/",
		Timeout:     2 * time.Second,
	}

	client, err := interservice_client.NewClientWithConfig(config)
	if err != nil {
		return false
	}

	makeurl := config.BaseURL + skustr
	body := map[string]string{
		// "hub_id": "",
		"skus": "",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return false
	}

	req := &http.Request{
		Url:     makeurl, // Use configured URL
		Body:    bytes.NewReader(bodyBytes),
		Timeout: 7 * time.Second,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}
	resp,_ := client.Get(req, "/")
	// if interSvcErr != nil {
	// 	log.Printf("Some InterSvcError: %+v", interSvcErr)
	// 	return !interSvcErr.StatusCode.Is4xx()
	// }
	if resp==nil{
		return false
	}
	return resp.IsSuccess()

	

}
func PushCreateOrderMessageToKafka(ctx context.Context, order *requests.Order) error {
	msg, err := json.Marshal(order.OrderItems)
	if err != nil {
		return fmt.Errorf("unable to marshal order items: %w", err)
	}

	headers := map[string]string{
		"event": "order_created",
	}

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

	ordermap := make(map[string]*requests.Order)
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
	for !Csv.IsEOF() {
		records, err := Csv.ReadNextBatch()
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV batch: %w", err)
		}

		// fmt.Println("Processing records:")
		// fmt.Println(records)
		for _, record := range records {
			orderNo := record[0]
			customerName := record[1]
			skuIDStr := record[2]
			hubIDStr := record[3]
			quantityStr := record[4]
			skuID, err := strconv.Atoi(skuIDStr)
			if err != nil {
				fmt.Println("invalid sku_id", skuIDStr, ":", err)
				continue
			}
			// hubID, err := strconv.Atoi(hubIDStr)
			// if err != nil {
			// 	fmt.Println("invalid sku_id", hubIDStr, ":", err)
			// 	continue
			// }

			if !validateSKU(skuID) {
				fmt.Println("/n")
				fmt.Println("/n")
				fmt.Println("/n")
				fmt.Println("sku -> ", skuIDStr, "doesn't exists")
				fmt.Println("/n")
				fmt.Println("/n")
				fmt.Println("/n")
				continue
			}

			quantity, err := strconv.Atoi(quantityStr)
			if err != nil {
				return nil, fmt.Errorf("invalid quantity %s: %v", quantityStr, err)
			}

			orderKey := fmt.Sprintf("%s-%s", orderNo, customerName)
			order, exists := ordermap[orderKey]
			if !exists {
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

			orderItem := requests.OrderItem{
				HubId:    hubIDStr,
				OrderID:  orderNo,
				SKUID:    skuIDStr,
				Quantity: quantity,
			}
			order.OrderItems = append(order.OrderItems, orderItem)
		}
	}
	var orders []*requests.Order
	for _, order := range ordermap {
		orders = append(orders, order)
	}

	fmt.Println("Final orders:")
	// for _, order := range orders {
	// 	fmt.Printf("Order No: %s, Customer: %s, Total Items: %d\n", order.OrderNo, order.CustomerName, len(order.OrderItems))
	// }

	return orders, nil
}
func ParseCSV(filePath string, ctx context.Context) {
	orders, err := ExtractFromCsv(filePath)
	if err != nil {
		fmt.Printf("\nFailed to parse CSV (%s): %v\n", filePath, err)
		return
	}

	for _, order := range orders {

		err := repository.CreateOrder(ctx, order)
		if err != nil {
			log.Fatalf("Could not create order %s: %v", order.OrderNo, err)
		}
		fmt.Println("/n")
		fmt.Println("/n")
		fmt.Println("/n")
		fmt.Println("huaa kya")
		fmt.Println("/n")
		fmt.Println("/n")
		if err := PushCreateOrderMessageToKafka(ctx, order); err != nil {
			fmt.Printf("Failed to publish order %s: %v\n", order.OrderNo, err)
		} else {
			fmt.Printf("Published order %s successfully\n", order.OrderNo)
		}
	}
}

func ConvertControllerReqToServiceReqParseCsv(ctx context.Context, CreateBulkCsv *requests.CSVUploadRequest) (string, error) {

	message := &sqs.Message{
		GroupId:       "csv-1665",
		Value:         []byte(CreateBulkCsv.FilePath),
		ReceiptHandle: "order-csv8",
		// DeduplicationId: "dedup-499",
		DelayDuration: 0,
		Headers: map[string]string{
			"customer_id": CreateBulkCsv.CustomerId,
		},
	}
	SendMessage(ctx, message)

	log.Println("Message sent successfully to SQS")
	return "", nil

}
