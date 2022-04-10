package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	const timeout = 2 * time.Second
	heartbeat, results := doWork(ctx, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			cancel()
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}
