package main

import (
	"sync"
	"time"
)

/*
	This util aims to help create recurring tasks
	with ease. Loops that implement the LoopResumable
	interface can pause and resume a given task.
	A loop implementation can also have an interval to
	add a delay between each task to reduce the performance
	impact of certain demanding tasks.
*/

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

// Creates a resumable loop that calls a callback. It can pause and resume, it can
// also wait for a period of time before calling the callback again.
func NewResumableLoop(interval time.Duration, cb func(LoopResumable)) *WorkLoop {
	return &WorkLoop{running: false, paused: false, interval: interval, cb: cb}
}

// Starts the work loop, this loop will
// run until the user pauses. After it is
// paused, in order to resume the Resume method
// calls startWork again to recreate the loop
func (w *WorkLoop) startWork() {
	if w.running == false {
		return
	}

	w.mu.Lock()
	w.paused = false
	w.mu.Unlock()

	for w.paused == false {
		w.cb(w) // Calls back the user code
		if w.interval > 0 {
			time.Sleep(w.interval)
		}
	}
}

/*
This implementation consists of two loops
the work loop and the block loop
the work loop is running until the user pauses
the block loop runs until the uses stops, until
then the job of the block loop is to prevent the
program from continuing if the Start method is called
with block=true

The use of two loops is used to avoid using an if statement
at every iteration of the loop to check if the loop
is paused eg: if w.paused; continue;
*/
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

// Recreates the work loop, because it was
// destroyed in order to pause
func (w *WorkLoop) Resume() {
	if w.paused && w.running {
		go w.startWork()
	}
}
