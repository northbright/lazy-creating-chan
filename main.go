package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	Name     string
	progress atomic.Value
	mu       sync.Mutex
	Done     chan struct{}
}

func NewTask(name string) *Task {
	return &Task{Name: name, Done: make(chan struct{})}
}

func (t *Task) ProgressCh() <-chan int {
	ch := t.progress.Load()
	if ch != nil {
		return ch.(chan int)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	ch = t.progress.Load()
	if ch == nil {
		ch = make(chan int)
		t.progress.Store(ch)
		fmt.Printf("channel created in ProgressCh()\n")
	}
	return ch.(chan int)
}

func (t *Task) Run() {
	fmt.Printf("Task: %v is running...\n", t.Name)

	for i := 0; i <= 100; i += 10 {
		ch, _ := t.progress.Load().(chan int)
		if ch != nil {
			ch <- i
		}
		time.Sleep(time.Millisecond * 50)
	}

	fmt.Printf("Task: %v goroutine exits\n", t.Name)
	t.Done <- struct{}{}
}

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
}
