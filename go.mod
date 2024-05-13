module github.com/koterin/broker/rabbitmq

go 1.18

replace github.com/koterin/broker/rabbitmq/pserver => ./pserver

require (
	github.com/google/uuid v1.6.0
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/rabbitmq/rabbitmq-stream-go-client v1.4.2
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
)
