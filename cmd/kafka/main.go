package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {

	brokers := []string{"172.22.0.1:29092"}
	topic := "test-topic"
	groupID := "test-consumer-group"

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	kafkaConnection, _ := kafka.Dial("tcp", "172.22.0.1:29092")

	err := kafkaConnection.CreateTopics(topicConfig)
	if err != nil {
		fmt.Println("Error creating the topic %w", err)
		return
	}

	readerConfig := kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	}

	// Set up Kafka writer configuration
	writerConfig := kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	}

	log.Println("Starting Kafka Consumer and Producer...")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	writer := kafka.NewWriter(writerConfig)
	reader := kafka.NewReader(readerConfig)

	defer kafkaConnection.Close()
	defer reader.Close()
	defer writer.Close()

	type myCustomMessage struct {
		Id    int    `json:"id"`
		Value string `json:"value"`
	}

	customMessage := myCustomMessage{
		Id:    1,
		Value: "myValue",
	}

	byteCustomMessage, err := json.Marshal(customMessage)

	if err != nil {
		log.Fatal("Error on json marshal:", err)
	}

	// write a message
	message := kafka.Message{
		Key:   []byte("my-key"),
		Value: byteCustomMessage,
	}

	err = writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Fatal("Error writing message:", err)
	}

	log.Println("Message produced to Kafka")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			fmt.Println("Termination signal received. Closing Kafka consumer and producer.")
			return
		default:
			// Read recent message from Kafka
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			log.Println("Received message")
			log.Printf("key: %v", string(msg.Key))

			var message myCustomMessage

			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("value: %v", message)
		}
	}

}
