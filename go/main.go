package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

func onMessage(msg *kafka.Message, err error) {
	if err == nil {
		// fmt.Println(string(msg.Value), rand.Intn(50))
	} else if err.(kafka.Error).IsFatal() {
		// The client will automatically try to recover from all errors.
		// Timeout is not considered an error because it is raised by
		// ReadMessage in absence of messages.
		fmt.Println("Consumer error:", err)
	}
}

func main() {
	cn := 0
	wk := NewResumableLoop(2*time.Second, func(lw LoopResumable) {
		cn += 1
		fmt.Println("working", cn)
	})

	defer wk.Stop()

	topic := "myTopic"

	// CONSUMER EXAMPLE
	c1, err := NewConsumerLoop(onMessage, "kafka", "myGroup", "c", []string{topic})

	// c1.SetMaxTasks(3)
	c1.SetTaskDelay(1 * time.Second)

	if err != nil {
		panic(err)
	}

	c1.Start(false)

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

		for i := 0; i < 20; i++ {
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte("very special message"),
			}, nil)
		}

		// Wait for message deliveries before shutting down
		// p.Flush(15 * 1000)

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/start", func(c *gin.Context) {
		wk.Start(false)
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.GET("/pause", func(c *gin.Context) {
		wk.Pause()
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.GET("/resume", func(c *gin.Context) {
		wk.Resume()
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.GET("/stop", func(c *gin.Context) {
		wk.Stop()
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
