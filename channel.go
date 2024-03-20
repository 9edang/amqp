package amqp

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/rabbitmq/amqp091-go"
)

type channel struct {
	*amqp091.Channel

	queue  amqp091.Queue
	closed int32

	mu *sync.Mutex
}

func (channel *channel) Close() error {
	if channel.IsClosed() {
		return nil
	}
	atomic.StoreInt32(&channel.closed, 1)
	return channel.Channel.Close()
}

func (channel *channel) IsClosed() bool {
	return atomic.LoadInt32(&channel.closed) == 1
}

func (channel *channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (err error) {
	channel.queue, err = channel.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, amqp091.Table(args))
	if err != nil {
		return err
	}
	return nil
}

func (channel *channel) Consume(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error) {
	channel.mu.Lock()
	defer channel.mu.Unlock()

	deloveries := make(chan Delivery)

	go func() {
		for {
			channel.mu.Lock()
			msgs, err := channel.Channel.ConsumeWithContext(ctx, queue, consumer, autoAck, exclusive, noLocal, noWait, amqp091.Table(args))
			channel.mu.Unlock()
			if err != nil {
				delay()
				continue
			}

			for msg := range msgs {
				deloveries <- deliveryWrap(msg)
			}

			delay()
			if channel.IsClosed() {
				close(deloveries)
				break
			}
		}

	}()

	return deloveries, nil
}

func (channel *channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error {
	return channel.Channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, amqp091.Table(args))
}

func (channel *channel) QueueBind(name, key, exchange string, noWait bool, args Table) error {
	return channel.Channel.QueueBind(name, key, exchange, noWait, amqp091.Table(args))
}

func (channel *channel) Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg Publishing) error {

	if channel.IsClosed() {
		return nil
	}

	return channel.Channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg.unWrap())
}
