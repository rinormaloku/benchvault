package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func measureAveragePerformance(client *api.Client, max int, steps int, currentMax int) (float64, int32) {
	var wg sync.WaitGroup
	totalMeasurements := max / (steps * 100)
	wg.Add(totalMeasurements)

	mx := sync.Mutex{}

	totalTime := 0 * time.Second
	responseTimeChan := make(chan time.Duration, totalMeasurements)

	// measure response time
	go func() {
		for {
			select {
			case timetaken := <-responseTimeChan:
				mx.Lock()
				totalTime += timetaken
				mx.Unlock()
				wg.Done()
			}
		}
	}()

	counter := atomic.Int32{}
	ctx := context.Background()

	// execute requests
	for i := 0; i < totalMeasurements; i++ {
		go func() {

			start := time.Now()
			path := fmt.Sprintf("ns%d/secret-%d", rand.Intn(currentMax+1), rand.Intn(max/steps))

			resp, err := client.KVv2("secret").Get(ctx, path)
			if err != nil {
				counter.Add(1)
				wg.Done()
				return
			}
			data := resp.Data["data"]
			_ = data
			//fmt.Printf("Secret read successfully: %v\n", resp.Data["data"])

			elapsed := time.Since(start)
			responseTimeChan <- elapsed
		}()
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
	}
	wg.Wait()
	microsecondsPerRequest := float64(int(totalTime.Microseconds())) / float64(totalMeasurements)
	errorcount := counter.Load()

	return microsecondsPerRequest, errorcount
}
