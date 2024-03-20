package amqp

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Publishing captures the client message sent to the server.  The fields
// outside of the Headers table included in this struct mirror the underlying
// fields in the content frame.  They use native types for convenience and
// efficiency.
type Publishing struct {
	// Application or exchange specific fields,
	// the headers exchange will inspect this field.
	Headers Table

	// Properties
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
	Priority        uint8     // 0 to 9
	CorrelationId   string    // correlation identifier
	ReplyTo         string    // address to to reply to (ex: RPC)
	Expiration      string    // message expiration spec
	MessageId       string    // message identifier
	Timestamp       time.Time // message timestamp
	Type            string    // message type name
	UserId          string    // creating user id - ex: "guest"
	AppId           string    // creating application id

	// The application specific payload of the message
	Body []byte
}

func (p Publishing) unWrap() amqp091.Publishing {
	return amqp091.Publishing{
		Headers:         amqp091.Table(p.Headers),
		ContentType:     p.ContentType,
		ContentEncoding: p.ContentEncoding,
		DeliveryMode:    p.DeliveryMode,
		Priority:        p.Priority,
		CorrelationId:   p.CorrelationId,
		ReplyTo:         p.ReplyTo,
		Expiration:      p.Expiration,
		MessageId:       p.MessageId,
		Timestamp:       p.Timestamp,
		Type:            p.Type,
		UserId:          p.UserId,
		AppId:           p.AppId,
		Body:            p.Body,
	}
}
