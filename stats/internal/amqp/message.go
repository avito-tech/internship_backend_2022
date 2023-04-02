package amqp

import (
	"time"
)

type KeyType string

const (
	TransactionEvent = "new-transaction"
)

type Message struct {
	Data      map[string]string `json:"data"`
	Timestamp time.Time         `json:"timestamp"`
	Key       string            `json:"key"`
}
