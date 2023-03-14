package amqp

import "time"

type KeyType string

const (
	TransactionEvent = "new-transaction"
)

type Message struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Key       string      `json:"key"`
}
