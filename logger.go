package amqp

type Logger interface {
	Printf(format string, args ...interface{})
}
