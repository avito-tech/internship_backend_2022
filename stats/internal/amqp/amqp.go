package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/radovsky1/internship_backend_2022/stats/cfg"
	"github.com/radovsky1/internship_backend_2022/stats/internal/domain"
	"github.com/radovsky1/internship_backend_2022/stats/internal/usecase"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type Consumer struct {
	amqpDial   *amqp.Connection
	amqpDialCh *amqp.Channel
	stop       chan bool
	config     config.RabbitMQ
}

func BuildConsumer(cfg *config.Config) (*Consumer, error) {
	amqpDial, err := amqp.Dial(cfg.RabbitMQ.DSN)
	if err != nil {
		return nil, err
	}

	ch, err := amqpDial.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		amqpDial:   amqpDial,
		stop:       make(chan bool),
		amqpDialCh: ch,
		config:     cfg.RabbitMQ,
	}, nil
}

func (c *Consumer) StartConsume(log *zap.Logger, uc usecase.ServiceUsecase) error {
	if err := c.amqpDialCh.Qos(1, 0, false); err != nil {
		return err
	}

	q, err := c.amqpDialCh.QueueDeclare(
		c.config.Queue, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	response, err := c.amqpDialCh.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for d := range response {
		var message Message

		if err := json.Unmarshal(d.Body, &message); err != nil {
			log.Error(fmt.Sprintf("err message format %+v", err))
			continue
		}

		ctx := context.Background()

		wg.Add(1)

		go func(msg Message) {
			defer wg.Done()

			log.Info(fmt.Sprintf("start consume message %+v", msg))

			switch KeyType(msg.Key) {
			case TransactionEvent:
				data := domain.Transaction{}

				for k, v := range msg.Data {
					switch k {
					case "ID":
						data.ID = uuid.MustParse(v)
					case "UserID":
						data.UserID = uuid.MustParse(v)
					case "ServiceID":
						data.ServiceID = uuid.MustParse(v)
					case "OrderID":
						data.OrderID = uuid.MustParse(v)
					case "Amount":
						data.Amount, _ = strconv.ParseFloat(v, 64)
					case "Status":
						data.Status = domain.TransactionStatus(v)
					case "CreatedAt":
						data.CreatedAt, _ = time.Parse(time.RFC3339, v)
					case "UpdatedAt":
						data.UpdatedAt, _ = time.Parse(time.RFC3339, v)
					}
				}

				err = uc.TransactionHandler(ctx, data)
			}

			if err != nil {
				log.Error(fmt.Sprintf("err consume message %+v", err))
				return
			}
		}(message)
	}

	wg.Wait()
	c.stop <- true

	return nil
}

func (c *Consumer) Close(ctx context.Context) error {
	if err := c.amqpDialCh.Close(); err != nil {
		return err
	}

	if err := c.amqpDial.Close(); err != nil {
		return err
	}

	for {
		select {
		case <-c.stop:
			return nil
		case <-ctx.Done():
			return fmt.Errorf("some missing handlers :(")
		}
	}
}
