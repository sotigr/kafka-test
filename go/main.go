package main

import (
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
)

func onMessage(msg *kafka.Message, err error) {
	if err == nil {
		fmt.Println(string(msg.Value))
	} else if !err.(kafka.Error).IsTimeout() {
		// The client will automatically try to recover from all errors.
		// Timeout is not considered an error because it is raised by
		// ReadMessage in absence of messages.
		fmt.Println("Consumer error:", err, msg)
	}
}

func main() {
	topic := "myTopic"

	// CONSUMER EXAMPLE
	c1, err := NewConsumerWorker(onMessage, "kafka", "myGroup", "c", topic)

	if err != nil {
		panic(err)
	}

	c1.Start()

	// PUBLISHER EXAMPLE
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error",
			})
			return
		}

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("very special message"),
		}, nil)
		// Wait for message deliveries before shutting down
		p.Flush(15 * 1000)

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
