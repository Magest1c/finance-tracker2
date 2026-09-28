package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan string, 1)

	go sendMessage(messages)

	fmt.Println("main is waiting")

	time.Sleep(2 * time.Second)

	fmt.Println("main receives message")

	message := <-messages
	fmt.Println("received:", message)
}

func sendMessage(message chan string) {
	fmt.Println("before sending")
	message <- "hello"
	
	fmt.Println("after sending")
}
