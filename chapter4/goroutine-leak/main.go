package main

import (
	"fmt"
	"math/rand"
	"time"
)

func goroutineCancel() {
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do someting
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}

func blockedGoroutineCancel() {
	newRandStream := func(done <-chan interface{}) (<-chan int, <-chan interface{}) {
		randStream, terminated := make(chan int), make(chan interface{})
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			defer close(terminated)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream, terminated
	}

	done := make(chan interface{})
	randStream, terminated := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	<-terminated
}

func main() {
	goroutineCancel()
	blockedGoroutineCancel()
}
