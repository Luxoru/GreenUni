package concurrency

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {

	pool := NewThreadPool(10, 10)

	pool.Start()

	pool.Submit(func() {
		Add()
	})

	time.Sleep(time.Second * 2)

}

func Add() {

	fmt.Println("1")
}
