package coroutine

import (
	"fmt"
	"testing"
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