package config

type Config struct {
	DatabaseURL          string `env:"DATABASE_URL,required"`
	Port                 string `env:"PORT,required"`
	ServiceName          string `env:"SERVICE_NAME,required"`
	OTLPExporterEndpoint string `env:"OTLP_EXPORTER_ENDPOINT,required"`
}
