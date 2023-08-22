package main

import (
	"log"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/streadway/amqp"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Consume
	msgs, err := ch.Consume( //env vars
		"queuename", // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// Connect to memcached
	mc := memcache.New(os.Getenv("MEMCACHE_URL"))

	// Loop over messages and write them to memcached
	for msg := range msgs {
		start := time.Now()

		// Write message body to memcached
		item := &memcache.Item{
			Key:        msg.MessageId, // job id
			Value:      msg.Body,      // json body
			Expiration: 0,             // ?
		}
		err = mc.Set(item)
		if err != nil {
			log.Printf("Failed to write message: %v", err)
			continue
		}

		// Print processing time for debugging
		log.Printf("Processed message %s in %v", msg.MessageId, time.Since(start))
	}
}
