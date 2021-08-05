package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Task can report progress when doing the work.
type Task struct {
	Name string
	// Used to store the progress channel.
	progress atomic.Value
	mu       sync.Mutex
	Done     chan struct{}
}

// NewTask creates the task by given task name.
func NewTask(name string) *Task {
	return &Task{Name: name, Done: make(chan struct{})}
}

// ProgressCh returns a channel to receive the task progress.
func (t *Task) ProgressCh() <-chan int {
	ch := t.progress.Load()
	if ch != nil {
		return ch.(chan int)
	}

	// Even load / save atomic.Value is goroutine safe,
	// still need mutex to protect the "transaction(load and store atomic.Value)" between differents goroutines.
	t.mu.Lock()
	defer t.mu.Unlock()
	ch = t.progress.Load()
	// Lazily created at first ProgressCh() is called.
	if ch == nil {
		ch = make(chan int)
		t.progress.Store(ch)
		fmt.Printf("channel created in ProgressCh()\n")
	}
	return ch.(chan int)
}

// Run starts the task work.
func (t *Task) Run() {
	fmt.Printf("Task: %v is running...\n", t.Name)

	for i := 0; i <= 100; i += 10 {
		// Send progress data to the channel if the channel exists.
		ch, _ := t.progress.Load().(chan int)
		if ch != nil {
			ch <- i
		}
		time.Sleep(time.Millisecond * 50)
	}

	fmt.Printf("Task: %v goroutine exits\n", t.Name)
	t.Done <- struct{}{}
}

// TestA is the most common use case.
// Starts the task in a new goroutine and create a for-select loop at once.
// Task's progress channel will be created at the first ProgressCh() is called(immediately after task goroutine started).
// It can capture the whole progress from 0 - 100.
func TestA() {
	t := NewTask("A")

	go func() {
		t.Run()
	}()

	for {
		select {
		case p := <-t.ProgressCh():
			fmt.Printf("Task: %v: %d%%\n", t.Name, p)
		case <-t.Done:
			fmt.Printf("Task: %v: done\n", t.Name)
			return
		}
	}
}

// TestB has a delay to create the for-select loop(also the progress channel),
// so it may lose some progress data(0 - 20) until the ProgressCh() is called.
func TestB() {
	t := NewTask("B")

	go func() {
		t.Run()
	}()

	time.Sleep(time.Millisecond * 200)

	for {
		select {
		case p := <-t.ProgressCh():
			fmt.Printf("Task: %v: %d%%\n", t.Name, p)
		case <-t.Done:
			fmt.Printf("Task: %v: done\n", t.Name)
			return
		}
	}
}

// TestC ignores the progress data and the progress channel will not be created at all.
func TestC() {
	t := NewTask("C")
	go func() {
		t.Run()
	}()

	<-t.Done
}

func main() {
	TestA()
	TestB()
	TestC()

	// Output:
	//channel created in ProgressCh()
	//Task: A is running...
	//Task: A: 0%
	//Task: A: 10%
	//Task: A: 20%
	//Task: A: 30%
	//Task: A: 40%
	//Task: A: 50%
	//Task: A: 60%
	//Task: A: 70%
	//Task: A: 80%
	//Task: A: 90%
	//Task: A: 100%
	//Task: A goroutine exits
	//Task: A: done
	//Task: B is running...
	//channel created in ProgressCh()
	//Task: B: 40%
	//Task: B: 50%
	//Task: B: 60%
	//Task: B: 70%
	//Task: B: 80%
	//Task: B: 90%
	//Task: B: 100%
	//Task: B goroutine exits
	//Task: B: done
	//Task: C is running...
	//Task: C goroutine exits
}
