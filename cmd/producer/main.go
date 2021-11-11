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

	// declare an exchange
	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare an exchange", err)
	}

	body := bodyFrom(os.Args)
	// publish a message to the queue
	err = ch.Publish("logs", "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatal("failed to publish a message", err)
	}

	log.Printf(" [x] Sent %s\n", body)
}

func bodyFrom(args []string) (s string) {
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}

	return
}
