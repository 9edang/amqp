package amqp

import "time"

var _defaultDelayReconnect = 3 * 1000 * time.Millisecond

func delay() {
	time.Sleep(_defaultDelayReconnect)
}

func setDelay(t time.Duration) {
	_defaultDelayReconnect = t
}
