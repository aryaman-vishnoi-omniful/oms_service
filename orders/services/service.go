package services

import (
	"context"
	"fmt"
	"oms_service/orders/requests"
	"os"
	"strconv"
	"time"

	// "oms_service/orders/listners"
	// "time"

	// "fmt"
	"log"

	// "github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/sqs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var NewProducer = &sqs.Publisher{}

// var Queue_instance =&sqs.Queue{}


func SetProducer(ctx context.Context, queue *sqs.Queue) {
	NewProducer = sqs.NewPublisher(queue)
	// ok := "papa"
	// jd, _ := json.Marshal(ok/)
	// newmessage:=NewProducer.Publish(ctx,)
	

	// fmt.Println(NewProducer)

	// Queue_instance=queue
	// go listners.StartConsume()
}
func SendMessage(ctx context.Context,message *sqs.Message){
	
	err := NewProducer.Publish(ctx,message)
	if err != nil {
		log.Fatal("did not publish", err)
	}

}
// type CSVUploadRequest struct {
// 	FilePath string
// 	UserID   string
// }
func ExtractFromCsv(filePath string) ([]*requests.Order, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Map to group items by order_no and customer_name
	ordermap:= make(map[string]*requests.Order)

	// Initialize the CSV reader (based on your previous implementation)
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
	err = Csv.InitializeReader(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}

	// Process the records and group them by order_no and customer_name
	for !Csv.IsEOF() {
		var records csv.Records
		records, err := Csv.ReadNextBatch()
		if err != nil {
			log.Fatal(err)
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

			// Check if the order group for this order_no and customer_name already exists
			orderKey := fmt.Sprintf("%s-%s", orderNo, customerName)
			order, exists := ordermap[orderKey]
			if !exists {
				// If order doesn't exist, create a new order
				now := primitive.NewDateTimeFromTime(time.Now())
				order = &requests.Order{
					ID: primitive.NewObjectID(),
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
				SKUID:    skuID,
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

func ParseCSV(filePath string) {
	orders, err := ExtractFromCsv(filePath)
	if err != nil {
		// c.JSON(500, gin.H{"error": err.Error()})
		fmt.Println("\nfailed to parse csv with path : ", filePath)
		return
	}
	fmt.Println(orders)

	// Store each order in MongoDB (optional, depending on the flow)
	// for _, order := range orders {
	// 	if err := storeOrder(order); err != nil {
	// 		fmt.Print("\nparseCSV : falied to save order")
	// 		return
	// 	}
	// }

}

func ConvertControllerReqToServiceReqParseCsv(ctx context.Context,CreateBulkCsv *requests.CSVUploadRequest)(string,error){

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
	return "",nil


}