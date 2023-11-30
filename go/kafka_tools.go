package main

import (
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type TaskLoop interface {
	Start(block bool)
	Stop()
	SetMaxTasks(max int)
	SetTaskDelay(delay time.Duration)
}

// A simple implementation of the Loop interface
// for a kafka consumer

type worker struct {
	running   bool
	c         *kafka.Consumer
	cb        func(*kafka.Message, error)
	mu        sync.Mutex
	maxTasks  int
	curTasks  int
	taskDelay time.Duration
}

func NewConsumerLoop(onMessage func(*kafka.Message, error), servers string, groupId string, instanceId string, topics []string) (TaskLoop, error) {
	var cw TaskLoop

	// Initialize the kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
		"group.instance.id": instanceId,
	})

	if err != nil {
		return cw, err
	}

	c.SubscribeTopics(topics, nil)

	cw = &worker{
		running:   false,
		c:         c,
		cb:        onMessage,
		maxTasks:  -1, // -1 means infinite tasks
		curTasks:  0,
		taskDelay: time.Duration(0),
	}

	return cw, err
}

func (w *worker) SetMaxTasks(max int) {
	w.maxTasks = max
}

func (w *worker) SetTaskDelay(delay time.Duration) {
	w.taskDelay = delay
}

func (w *worker) runLoop() {
	for w.running {
		if w.maxTasks > 0 {
			if w.curTasks >= w.maxTasks {
				time.Sleep(w.taskDelay)
				continue
			}
			msg, err := w.c.ReadMessage(1 * time.Second)

			w.mu.Lock()
			w.curTasks += 1
			w.mu.Unlock()
			go func(msg *kafka.Message, err error) {

				w.cb(msg, err)

				w.mu.Lock()
				w.curTasks -= 1
				w.mu.Unlock()

			}(msg, err)
		} else {
			msg, err := w.c.ReadMessage(1 * time.Second)
			w.cb(msg, err)
			time.Sleep(w.taskDelay)
		}

	}
	w.c.Close()
}

func (w *worker) Start(block bool) {
	if w.running {
		return
	}
	w.mu.Lock()
	w.running = true
	w.mu.Unlock()

	// Read messages and give them to the user to handle
	// if block, block the thread otherwise use a goroutine
	if block {
		w.runLoop()
	} else {
		go func() {
			w.runLoop()
		}()
	}

}

func (w *worker) Stop() {
	w.mu.Lock()
	w.running = false
	w.mu.Unlock()
}
