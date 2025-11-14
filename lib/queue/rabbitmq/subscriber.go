package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type Handler func(ctx context.Context, body []byte) error

func NewSubscriber(url string) (*Subscriber, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Subscriber{
		conn:    conn,
		channel: ch,
	}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, queue string, handler Handler) error {
	q, err := s.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := s.channel.Consume(
		q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				if err := handler(ctx, msg.Body); err != nil {
					msg.Nack(false, true) // requeue on error
					continue
				}

				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (s *Subscriber) Close() error {
	if err := s.channel.Close(); err != nil {
		return err
	}
	return s.conn.Close()
}
