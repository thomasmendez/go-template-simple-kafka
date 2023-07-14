package main

import (
	"context"
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

	// Set up Kafka reader configuration
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

	fmt.Println("Starting Kafka Consumer and Producer...")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	writer := kafka.NewWriter(writerConfig)
	reader := kafka.NewReader(readerConfig)

	defer reader.Close()
	defer writer.Close()

	// write a message
	message := kafka.Message{
		Key:   []byte("key"),
		Value: []byte("value"),
	}

	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Println("Error writing message:", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// Start consuming and producing messages
	for {
		select {
		case <-signals:
			fmt.Println("Termination signal received. Closing Kafka consumer and producer.")
			return
		default:
			// Read a message from Kafka
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatal("Error reading message from Kafka:", err)
			}

			// Process the message
			fmt.Printf("Received message: %s\n", string(msg.Value))

			// Modify the message and produce it back to Kafka
			newMessage := kafka.Message{
				Key:   []byte("modified-key"),
				Value: []byte("Modified: " + string(msg.Value)),
			}

			err = writer.WriteMessages(context.Background(), newMessage)
			if err != nil {
				log.Fatal("Error writing message to Kafka:", err)
			}

			fmt.Println("Message produced to Kafka")
		}
	}

	// serverErrors := make(chan error, 1)
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// select {
	// case err := <-serverErrors:

	// 	fmt.Println("server error: %w", err)

	// case sig := <-shutdown:
	// 	fmt.Println("shutdown started with sig: %w", sig)

	// 	// Give outstanding requests a deadline for completion.
	// 	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	// 	ctx.Done()
	// 	defer cancel()

	// 	// Asking listener to shutdown and shed load
	// 	// if err := api.Shutdown(ctx); err != nil {
	// 	// 	return fmt.Errorf("could not stop server gracefully: %w", err)
	// 	// }

	// }

	// return

}
