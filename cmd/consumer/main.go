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

	// declare an exchange
	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare an exchange", err)
	}

	// Firstly, whenever we connect to Rabbit we need a fresh, empty queue.
	// Secondly, once we disconnect the consumer the queue should be automatically deleted.
	// When we supply queue name as an empty string, we create a non-durable queue with a generated name
	// When the connection that declared it closes, the queue will be deleted because it is declared as exclusive.
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatal("failed to declare a queue", err)
	}

	// Now we need to tell the exchange to send messages to our queue.
	// That relationship between exchange and a queue is called a binding.
	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	if err != nil {
		log.Fatal("failed to bind a queue", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to register a consumer", err)
	}

	// since it will push us messages asynchronously,
	// we will read the messages from a channel
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
