package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.SetPrefix("[LOG] ")
	log.SetFlags(3)

	log.Printf("Producer server started successfully")
	rmqServerURL := os.Getenv("RMQ_SERVER_URL")
	log.Printf("RMQ server URL is: %s", rmqServerURL)
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(rmqServerURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// https://www.rabbitmq.com/tutorials/tutorial-one-go.html
	q := queueDeclare("notesDur", ch) // declare a queue, then excange will publish to this queue.

	log.Println("Enter your message: ")

	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("An error occured while reading input. Please try again", err)
			return
		}

		input = strings.TrimSuffix(input, "\n")

		if strings.Compare(input, "exit") == 0 {
			log.Println("Gracefully shutting down...")
			break
		}
		publishMsg(ch, q.Name, input)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func publishMsg(ch *amqp.Channel, name string, input string) {
	err := ch.Publish(
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(input),
		})
	failOnError(err, "Failed to publish a message")
}

func queueDeclare(name string, ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return q
}
