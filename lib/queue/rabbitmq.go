package queue

import (
	"blink/lib/telemetry"
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RabbitMQPubSub struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tracer  trace.Tracer
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
		tracer:  otel.Tracer("rabbitmq"),
	}, nil
}

func (r *RabbitMQPubSub) Subscribe(ctx context.Context, queue string, handler Handler) error {
	ctx, span := r.tracer.Start(ctx, "rabbitmq.RabbitMQPubSub.subscribe")
	defer span.End()

	telemetry.SetAttributes(
		ctx,
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.destination", queue),
		attribute.String("messaging.operation", "receive"),
	)

	q, err := r.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		telemetry.RecordError(ctx, err)
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
		telemetry.RecordError(ctx, err)
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

				msgCtx := r.extractTraceContext(context.Background(), msg.Headers)
				msgCtx, msgSpan := r.tracer.Start(msgCtx, "rabbitmq.RabbitMQPubSub.process",
					trace.WithSpanKind(trace.SpanKindConsumer),
					trace.WithAttributes(
						attribute.String("messaging.system", "rabbitmq"),
						attribute.String("messaging.destination", queue),
						attribute.Int("messaging.message_payload_size_bytes", len(msg.Body)),
					),
				)

				if err := handler(msgCtx, msg.Body); err != nil {
					telemetry.RecordError(ctx, err)
					msg.Nack(false, true)
					msgSpan.End()
					continue
				}

				msgSpan.SetStatus(codes.Ok, "message processed successfully")
				msgSpan.End()
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQPubSub) Publish(ctx context.Context, queue string, body []byte) error {
	ctx, span := r.tracer.Start(ctx, "rabbitmq.RabbitMQPubSub.publish",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "rabbitmq"),
			attribute.String("messaging.destination", queue),
			attribute.Int("messaging.message_payload_size_bytes", len(body)),
		),
	)
	defer span.End()

	q, err := r.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		telemetry.RecordError(ctx, err)
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	headers := r.injectTraceContext(ctx)

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
			Headers:      headers,
		},
	)
	if err != nil {
		telemetry.RecordError(ctx, err)
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

func (r *RabbitMQPubSub) injectTraceContext(ctx context.Context) amqp.Table {
	headers := make(amqp.Table)

	propagator := otel.GetTextMapPropagator()

	propagator.Inject(ctx, &amqpHeaderCarrier{headers: headers})

	return headers
}

func (r *RabbitMQPubSub) extractTraceContext(ctx context.Context, headers amqp.Table) context.Context {
	propagator := otel.GetTextMapPropagator()

	return propagator.Extract(ctx, &amqpHeaderCarrier{headers: headers})
}

type amqpHeaderCarrier struct {
	headers amqp.Table
}

func (c *amqpHeaderCarrier) Get(key string) string {
	if val, ok := c.headers[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

func (c *amqpHeaderCarrier) Set(key, value string) {
	c.headers[key] = value
}

func (c *amqpHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.headers))
	for k := range c.headers {
		keys = append(keys, k)
	}
	return keys
}
