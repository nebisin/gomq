package main

import (
	"log"
	"os"
	"strings"

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
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare a queue", err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatal("failed to set QoS", err)
	}

	body := bodyFrom(os.Args)
	// publish a message to the queue
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	if err != nil {
		log.Fatal("failed to publish a message", err)
	}
}

func bodyFrom(args []string) (s string) {
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}

	return
}
