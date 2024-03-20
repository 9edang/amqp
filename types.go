package amqp

import (
	"context"
	"time"
)

type Table map[string]interface{}

type Connection interface {
	SetLogger(l Logger)
	SetDelay(t time.Duration)
	Channel() (Channel, error)
	IsClosed() bool
	Close() error
	Dial() error
}

type Channel interface {
	Close() error
	IsClosed() bool
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (err error)
	QueueBind(name, key, exchange string, noWait bool, args Table) error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error
	Consume(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error)
	Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg Publishing) error
}
