package listners

import (
	"context"
	"fmt"
	"oms_service/orders/services"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	// "github.com/joho/godotenv"
	"log"
	// "os"
	"time"
)

func StartConsume(queuestring string,ctx context.Context) {

	
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	sqsClient := sqs.NewFromConfig(cfg)
	for {
		recieveMessages := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queuestring),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     1,
		}
		resp, err := sqsClient.ReceiveMessage(ctx, recieveMessages)
		if err != nil {
			log.Fatal(err)
		}
		for _, message := range resp.Messages {
			// fmt.Println(*message.Body)
			fmt.Println("gjnjdmnjgvfkmdgnjvfkmdmdecfmdcfdnfvhjdfccfm")
			fmt.Println(*message.Body)
			services.ParseCSV(*message.Body,ctx)

			_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queuestring),
				ReceiptHandle: message.ReceiptHandle,
			})

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Message deleted")
			}
		}
		time.Sleep(5 * time.Second)

	}

}