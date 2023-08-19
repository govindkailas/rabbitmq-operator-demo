module github.com/koterin/broker/rabbitmq

go 1.18

replace github.com/koterin/broker/rabbitmq/pserver => ./pserver

require github.com/rabbitmq/amqp091-go v1.3.4
