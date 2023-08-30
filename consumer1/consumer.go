package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	log.SetPrefix("[LOG] ")
	log.SetFlags(3)

	log.Printf("Consumer server started successfully")

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	rmqServerURL := os.Getenv("RMQ_SERVER_URL")
	log.Printf("RMQ server URL is: %s", rmqServerURL)
	conn, err := amqp.Dial(rmqServerURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rabbitmq-demo", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		2,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// Postgres DB connection settings
	// export POSTGRES_URL="postgres://user:password@localhost:5432/mydb"
	db, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close(context.Background())

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Println("Received a message: ", string(d.Body))

			_, err = db.Exec(context.Background(), "INSERT INTO consumer_data (data,consumer_name) VALUES ($1,$2)", string(d.Body), "consumer1")
			if err != nil {
				log.Fatalf("Failed to insert data into PostgreSQL: %v", err)
			}
			d.Ack(false) // we are good to mark the msg as acked now as its inserted to the DB
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
