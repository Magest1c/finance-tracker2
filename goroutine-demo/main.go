package main

import (
	"context"
	"fmt"
	"time"
)

func waitForWork(ctx context.Context) {

	fmt.Println("work started")

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("work completed")

	case <-ctx.Done():
		fmt.Println("work cancelled:", ctx.Err())
	}
}

func main() {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	waitForWork(ctx)

	fmt.Println("main finished")
}
