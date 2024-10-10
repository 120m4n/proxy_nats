package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

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

    // Publish a message every 7 seconds
    ticker := time.NewTicker(7 * time.Second)
    defer ticker.Stop()

    go func() {
        for t := range ticker.C {
            message := fmt.Sprintf("Current time: %s", t.UTC().Format(time.RFC3339))
            if err := nc.Publish("example.topic", []byte(message)); err != nil {
                log.Println("Error publishing message:", err)
            } else {
                fmt.Println("Published message:", message)
            }
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the publisher
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    fmt.Println("Shutting down...")
}