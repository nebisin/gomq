package main

import (
	"log"

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
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare a queue", err)
	}

	// publish a message to the queue
	body := "Hello World!"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatal("failed to publish a message", err)
	}
}
