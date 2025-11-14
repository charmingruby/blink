package config

type Config struct {
	Port                 string `env:"PORT,required"`
	ServiceName          string `env:"SERVICE_NAME,required"`
	OTLPExporterEndpoint string `env:"OTLP_EXPORTER_ENDPOINT,required"`
}
