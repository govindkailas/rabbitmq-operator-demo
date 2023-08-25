package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("%s ", err)
		os.Exit(1)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	//stream.SetLevelInfo(logs.DEBUG)
	fmt.Println("Configure a load-balancer TLS for RabbitMQ")
	fmt.Println("Connecting to RabbitMQ streaming ...")

	// load balancer address in TLS
	addressResolver := stream.AddressResolver{
		Host: os.Getenv("rmqLBIP"),
		Port: 5552,
	}
	conf := &tls.Config{InsecureSkipVerify: true}

	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(addressResolver.Host).
			SetPort(addressResolver.Port).
			IsTLS(false).
			SetTLSConfig(conf).
			SetAddressResolver(addressResolver).
			SetUser(os.Getenv("username")).
			SetPassword(os.Getenv("password")).
			SetMaxProducersPerClient(5))

	CheckErr(err)

	/// We create a few streams, in order to distribute the streams across the cluster
	var streamsName []string
	for i := 0; i < 3; i++ {
		streamsName = append(streamsName, uuid.New().String())
	}

	for _, streamName := range streamsName {
		fmt.Printf("Create stream %s ...\n", streamName)
		err = env.DeclareStream(streamName,
			&stream.StreamOptions{
				MaxLengthBytes: stream.ByteCapacity{}.GB(2),
			},
		)
	}

	CheckErr(err)
	var producers []*stream.Producer
	// The producer MUST connect to the leader stream
	// here the AddressResolver try to get the leader
	// if fails retry
	for _, streamName := range streamsName {
		fmt.Printf("Create producer for %s ...\n", streamName)
		producer, err := env.NewProducer(streamName, nil)
		producers = append(producers, producer)
		CheckErr(err)
	}

	// just publish some message
	for i := 0; i < 50; i++ {
		for _, producer := range producers {
			err := producer.Send(amqp.NewMessage([]byte("hello_world_" + strconv.Itoa(i))))
			CheckErr(err)
		}
	}

	handleMessages := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		fmt.Printf("consumer name: %s, text: %s \n ", consumerContext.Consumer.GetName(), message.Data)
	}

	// the consumer can connect to the leader o follower
	// the AddressResolver just resolve the ip
	var consumers []*stream.Consumer
	for _, streamName := range streamsName {
		fmt.Printf("Create consumer for %s ...\n", streamName)
		consumer, err := env.NewConsumer(
			streamName,
			handleMessages,
			stream.NewConsumerOptions().
				SetConsumerName(uuid.New().String()).            // set a random name
				SetOffset(stream.OffsetSpecification{}.First())) // start consuming from the beginning
		CheckErr(err)
		consumers = append(consumers, consumer)
	}

	/// check on the UI http://localhost:15673/#/stream/connections
	// the producers are connected to the leader node
	/// the consumers random nodes it doesn't matter

	fmt.Println("Hit Enter to stop ")
	_, _ = reader.ReadString('\n')
	for _, streamName := range streamsName {
		fmt.Printf("Delete stream %s ...\n", streamName)
		err = env.DeleteStream(streamName)
	}
	CheckErr(err)
	err = env.Close()
	CheckErr(err)

}
