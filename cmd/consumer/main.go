package main

import (
	"bytes"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect the RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal("error while dialing the rabbit mq", err)
	}
	defer conn.Close()

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("failed to open a channel", err)
	}
	defer ch.Close()

	// declare a queue
	// because we might start the consumer before the publisher
	// we want to make sure the queue exist
	q, err := ch.QueueDeclare("task_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare a queue", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to register a consumer", err)
	}

	// since it will push us messages asynchronously,
	// we will read the messages from a channel
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Println("Done")
			d.Ack(false)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
