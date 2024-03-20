package amqp

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type connection struct {
	*amqp091.Connection
	cfg    amqp091.Config
	dsn    string
	logger Logger

	mu *sync.Mutex
}

func New(dsn string) Connection {
	return new(dsn, nil)
}

func NewTLS(dsn string, tls *tls.Config) Connection {
	return new(dsn, tls)
}

func new(dsn string, tls *tls.Config) *connection {
	cfg := amqp091.Config{
		ChannelMax:      defaultChannelMax,
		Heartbeat:       defaultHeartbeat,
		Locale:          defaultLocale,
		TLSClientConfig: tls,
	}
	return &connection{
		cfg:    cfg,
		dsn:    dsn,
		logger: logger,
		mu:     &sync.Mutex{},
	}
}
func (c *connection) Close() error {
	if c.Connection == nil {
		return nil
	}

	if c.Connection.IsClosed() {
		return nil
	}

	return c.Connection.Close()
}

func (c *connection) SetLogger(l Logger) {
	c.logger = l
}

func (c *connection) SetDelay(t time.Duration) {
	setDelay(t)
}

func (c *connection) Dial() error {
	var err error
	c.Connection, err = amqp091.DialConfig(c.dsn, c.cfg)
	if err != nil {
		return err
	}

	c.logger.Printf("[AMQP] connected to server")

	go c.dialReconnect()

	return nil
}

func (c *connection) dialReconnect() {
	for {
		notify, ok := <-c.Connection.NotifyClose(make(chan *amqp091.Error))
		if !ok {
			break
		}

		c.logger.Printf("[AMQP] server disconnect reason: %s", notify.Reason)

		for {
			delay()

			conn, err := amqp091.DialConfig(c.dsn, c.cfg)
			if err != nil {
				continue
			}

			c.mu.Lock()
			c.Connection = conn
			c.mu.Unlock()

			c.logger.Printf("[AMQP] connected to server")
			break
		}
	}
}

func (c *connection) Channel() (Channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &channel{
		Channel: ch,
		mu:      &sync.Mutex{},
	}

	go channel.channelReconnect(c)

	return channel, nil
}

func (ch *channel) channelReconnect(c *connection) {
	for {
		_, ok := <-ch.Channel.NotifyClose(make(chan *amqp091.Error))
		if !ok {
			break
		}

		for {
			delay()

			if c.Connection.IsClosed() {
				continue
			}

			newCh, err := c.Connection.Channel()
			if err != nil {
				continue
			}

			ch.mu.Lock()
			ch.Channel = newCh
			ch.mu.Unlock()

			break
		}
	}
}
