package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/nats-io/nats.go"
)

func main() {
    // Connect to NATS server
    natsUrl := os.Getenv("NATS_URL")
    if natsUrl == "" {
        natsUrl = nats.DefaultURL
    }
	
    nc, err := nats.Connect(natsUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    // Subscribe to a topic
    sub, err := nc.Subscribe("example.topic", func(msg *nats.Msg) {
        fmt.Printf("Received message: %s\n", string(msg.Data))
    })
    if err != nil {
        log.Fatal(err)
    }
    defer sub.Unsubscribe()

    // Wait for interrupt signal to gracefully shutdown the subscription
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    fmt.Println("Shutting down...")
}