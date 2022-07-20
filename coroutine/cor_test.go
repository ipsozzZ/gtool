package coroutine

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestCor(t *testing.T) {

}

func leak() {
	ch := make(chan int32)

	go func() {
		val := <-ch
		fmt.Println("I received a value:", val)
	}()
}

// context控制go
func TestGoContext(t *testing.T) {
	tr := NewTracker()
	go tr.Run()

	_ = tr.Event(context.Background(), "test1")
	_ = tr.Event(context.Background(), "test2")
	_ = tr.Event(context.Background(), "test3")
	_ = tr.Event(context.Background(), "test4")
	_ = tr.Event(context.Background(), "test5")
	_ = tr.Event(context.Background(), "test6")
	_ = tr.Event(context.Background(), "test7")
	//time.Sleep(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	tr.Shutdown(ctx)
}

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{
		ch:   make(chan string, 10),
		stop: make(chan struct{}, 1),
	}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	fmt.Println("stop")
	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
		fmt.Println("stop")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}

// WaitGroup控制go
func TestGoSync(t *testing.T) {
	t1 := Tracker1{wg: sync.WaitGroup{}}
	t1.Event("test1")
	t1.Event("test2")
	t1.Event("test3")
	t1.Event("test4")
	t1.Event("test5")
	t1.Event("test6")

	t1.Shutdown()
}

type Tracker1 struct {
	wg sync.WaitGroup
}

func (t1 *Tracker1) Event(data string) {
	t1.wg.Add(1)

	go func() {
		defer t1.wg.Done()
		time.Sleep(1 * time.Millisecond)
		log.Println(data)
	}()
}

func (t1 *Tracker1) Shutdown() {
	t1.wg.Wait()
}
