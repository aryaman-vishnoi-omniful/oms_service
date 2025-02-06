package orders

import (
	"context"
	// "time"

	// "fmt"
	"log"

	"github.com/omniful/go_commons/sqs"
)

var NewProducer = &sqs.Publisher{}

// var Queue_instance =&sqs.Queue{}

func SetProducer(ctx context.Context, queue *sqs.Queue) {
	NewProducer = sqs.NewPublisher(queue)
	// ok := "papa"
	// jd, _ := json.Marshal(ok/)
	// newmessage:=NewProducer.Publish(ctx,)
	newmessage := &sqs.Message{
		GroupId:         "group-123",
		Value:           []byte("sample message"),
		ReceiptHandle:   "receipt-abc",
		Attributes:      map[string]string{"key1": "value1", "key2": "value2"},
		DeduplicationId: "dedup-456",
	}
	err := NewProducer.Publish(ctx, newmessage)
	if err != nil {
		log.Fatal("ni aaya", err)
	}

	// fmt.Println(NewProducer)

	// Queue_instance=queue
}
