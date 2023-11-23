package main

import (
	"sync"
	"time"
)

type WorkLoop interface {
	startWork()
	Start(block bool)
	Stop()
	Resume()
}

type LoopWorker struct {
	running  bool
	paused   bool
	interval time.Duration
	cb       func(*LoopWorker)
	mu       sync.Mutex
}

func NewLoopWorker(interval time.Duration, cb func(*LoopWorker)) *LoopWorker {
	return &LoopWorker{running: false, paused: false, interval: interval, cb: cb}
}

func (w *LoopWorker) startWork() {
	if w.running == false {
		return
	}

	w.mu.Lock()
	w.paused = false
	w.mu.Unlock()

	for w.paused == false {
		w.cb(w)
		if w.interval > 0 {
			time.Sleep(w.interval)
		}
	}
}

func (w *LoopWorker) Start(block bool) {
	if w.running {
		return
	}

	w.mu.Lock()
	w.running = true
	w.mu.Unlock()

	if block {
		w.startWork()
		for w.running {
		}
	} else {
		go func() {
			w.startWork()
			for w.running {
			}
		}()
	}

}

func (w *LoopWorker) Stop() {
	w.mu.Lock()
	w.paused = true
	w.running = false
	w.mu.Unlock()
}

func (w *LoopWorker) Pause() {
	w.mu.Lock()
	w.paused = true
	w.mu.Unlock()
}

func (w *LoopWorker) Resume() {
	if w.paused && w.running {
		go w.startWork()
	}
}
