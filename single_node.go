package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func SingleNode(ctx context.Context, toCall string) []byte {
	responseChannel := make(chan *Response, *totalCalls*2)

	benchTime := NewTimer()
	benchTime.Reset()

	wg := &sync.WaitGroup{}
	// TODO duration mode

	// 判断执行模式
	if *duration != "" {
		// duration 模式
		dur, err := time.ParseDuration(*duration)
		if err != nil {
			log.Fatalf("Invalid duration: %v", err)
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), dur)
		defer cancel()
	} else {
		// totalCalls 模式：无限超时
		ctx = context.Background()
	}

	for i := 0; i < *numConnections; i++ {
		go StartClient(
			ctx,
			toCall,
			*headers,
			*requestBody,
			*method,
			*disableKeepAlives,
			responseChannel,
			wg,
			*totalCalls,
		)
		wg.Add(1)
	}

	wg.Wait()

	result := CalcStats(
		responseChannel,
		benchTime.Duration(),
	)
	return result
}
