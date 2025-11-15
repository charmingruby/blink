package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPubSub struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type Handler func(ctx context.Context, body []byte) error

func NewRabbitMQPubSub(url string) (*RabbitMQPubSub, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQPubSub{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQPubSub) Subscribe(ctx context.Context, queue string, handler Handler) error {
	q, err := r.channel.QueueDeclare(
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

	msgs, err := r.channel.Consume(
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

func (r *RabbitMQPubSub) Publish(ctx context.Context, queue string, body []byte) error {
	q, err := r.channel.QueueDeclare(
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

	err = r.channel.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/protobuf",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (r *RabbitMQPubSub) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}

	return r.conn.Close()
}
