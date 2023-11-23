package main

import (
	"sync"
	"time"
)

type Loop interface {
	Start(block bool)
	Stop()
}

type LoopResumable interface {
	Start(block bool)
	Stop()
	Pause()
	Resume()
}

type WorkLoop struct {
	running  bool
	paused   bool
	interval time.Duration
	cb       func(LoopResumable)
	mu       sync.Mutex
}

func NewLoopWorker(interval time.Duration, cb func(LoopResumable)) *WorkLoop {
	return &WorkLoop{running: false, paused: false, interval: interval, cb: cb}
}

func (w *WorkLoop) startWork() {
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

func (w *WorkLoop) Start(block bool) {
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

func (w *WorkLoop) Stop() {
	w.mu.Lock()
	w.paused = true
	w.running = false
	w.mu.Unlock()
}

func (w *WorkLoop) Pause() {
	w.mu.Lock()
	w.paused = true
	w.mu.Unlock()
}

func (w *WorkLoop) Resume() {
	if w.paused && w.running {
		go w.startWork()
	}
}
