package config

type Config struct {
	SeverAddress         string `env:"SERVER_ADDRESS,required"`
	ServiceName          string `env:"SERVICE_NAME,required"`
	OTLPExporterEndpoint string `env:"OTLP_EXPORTER_ENDPOINT,required"`
	DatabaseURL          string `env:"DATABASE_URL,required"`
}
