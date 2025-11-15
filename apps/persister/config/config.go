package config

type Config struct {
	ServiceName          string `env:"SERVICE_NAME,required"`
	OTLPExporterEndpoint string `env:"OTLP_EXPORTER_ENDPOINT,required"`
	DatabaseURL          string `env:"DATABASE_URL,required"`
	QueueURL             string `env:"QUEUE_URL,required"`
	QueueName            string `env:"QUEUE_NAME,required"`
}
