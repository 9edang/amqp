package amqp

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	defaultChannelMax = (2 << 10) - 1
	defaultHeartbeat  = 10 * time.Second
	defaultLocale     = "en_US"
	logger            = logrus.New()
)
