module github.com/koterin/broker/rabbitmq

go 1.18

replace github.com/koterin/broker/rabbitmq/pserver => ./pserver

require github.com/rabbitmq/amqp091-go v1.3.4

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rabbitmq/rabbitmq-stream-go-client v1.2.0 // indirect
)
