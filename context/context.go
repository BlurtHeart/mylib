// implement standard library's context

package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func sleepAndTalk(ctx context.Context, duration time.Duration, data string) {
	select {
	case <-time.After(duration):
		fmt.Println(data)
	case <-ctx.Done():
		log.Print(ctx.Err())
	}
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*6)
	defer cancel()

	sleepAndTalk(ctx, time.Second*5, "hello")
}
