package main

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type worker struct {
	running bool
	c       *kafka.Consumer
	cb      func(*kafka.Message, error)
}

func NewConsumerWorker(onMessage func(*kafka.Message, error), servers string, groupId string, instanceId string, topic string) (Loop, error) {
	var cw Loop

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
		"group.instance.id": instanceId,
	})

	if err != nil {
		return cw, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	cw = worker{running: false, c: c, cb: onMessage}

	return cw, err
}

func (w worker) Start(block bool) {
	if w.running {
		return
	}
	w.running = true
	if block {
		for w.running {
			msg, err := w.c.ReadMessage(time.Second)
			w.cb(msg, err)
		}
		w.c.Close()
	} else {
		go func() {
			for w.running {
				msg, err := w.c.ReadMessage(time.Second)
				w.cb(msg, err)
			}
			w.c.Close()
		}()
	}

}

func (w worker) Stop() {
	w.running = false
}
