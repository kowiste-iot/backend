package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/nats-io/nats.go"
)

type WireMessage struct {
    ID        string
    Topic     string
    Data      string
    Timestamp time.Time
    Event     string
}

func prettyPrintJSON(msg *nats.Msg) string {
    var wireMsg WireMessage
    if err := json.Unmarshal(msg.Data, &wireMsg); err != nil {
        return string(msg.Data)
    }

    // Decode base64 data
    decodedData, err := base64.StdEncoding.DecodeString(wireMsg.Data)
    if err != nil {
        decodedData = []byte(wireMsg.Data)
    }

    // Try to parse decoded data as JSON
    var parsedData interface{}
    if err := json.Unmarshal(decodedData, &parsedData); err != nil {
        parsedData = string(decodedData)
    }

    // Create display struct with decoded data
    displayMsg := struct {
        ID        string      `json:"id"`
        Topic     string      `json:"topic"`
        Data      interface{} `json:"data"`
        Timestamp time.Time   `json:"timestamp"`
        Event     string      `json:"event"`
    }{
        ID:        wireMsg.ID,
        Topic:     wireMsg.Topic,
        Data:      parsedData,
        Timestamp: wireMsg.Timestamp,
        Event:     wireMsg.Event,
    }

    var prettyJSON bytes.Buffer
    encoder := json.NewEncoder(&prettyJSON)
    encoder.SetIndent("", "    ")
    if err := encoder.Encode(displayMsg); err != nil {
        return string(msg.Data)
    }
    return prettyJSON.String()
}

func StartSubscriber(url, subject string) {
    log.Printf("Connecting to NATS server: %s\n", url)
    log.Printf("use '>' for multi-level wildcard , '.*'  single-level wildcard : %s\n", url)
    log.Printf("Subscribing to subject: %s\n", subject)

    nc, err := nats.Connect(url)
    if err != nil {
        log.Fatal("Connection error:", err)
    }
    defer nc.Close()

    _, err = nc.Subscribe(subject, func(msg *nats.Msg) {
        fmt.Printf("\nReceived message:\n")
        fmt.Printf("Subject: %s\n", msg.Subject)
        fmt.Printf("Data:\n%s", prettyPrintJSON(msg))
        fmt.Println("-------------------")
    })
    if err != nil {
        log.Fatal("Subscription error:", err)
    }

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
}

func main() {
    url := flag.String("url", "nats://localhost:4222", "NATS server URL")
    subject := flag.String("topic", ">", "Subject to subscribe to")
    flag.Parse()

    StartSubscriber(*url, *subject)
}