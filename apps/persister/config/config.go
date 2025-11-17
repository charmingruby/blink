package config

type Config struct {
	ServiceName          string `env:"SERVICE_NAME,required"`
	OTLPExporterEndpoint string `env:"OTLP_EXPORTER_ENDPOINT,required"`
	PostgresURL          string `env:"POSTGRES_URL,required"`
	RabbitMQURL          string `env:"RABBITMQ_URL,required"`
	QueueName            string `env:"QUEUE_NAME,required"`
	RedisURL             string `env:"REDIS_URL,required"`
}
