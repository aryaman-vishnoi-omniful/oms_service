package kafka_producer

import (
	"github.com/omniful/go_commons/kafka"
)

type Producer struct {
	*kafka.ProducerClient
}

var clientInstance *Producer

func Get() *Producer {
	return clientInstance
}

func Set(client *kafka.ProducerClient) {
	clientInstance = &Producer{client}
}
