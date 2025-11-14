package config

type Config struct {
	Port                    string `env:"PORT,required"`
	ServiceName             string `env:"SERVICE_NAME,required"`
	OTLPExporterEndpoint    string `env:"OTLP_EXPORTER_ENDPOINT,required"`
	RecallGRPCServerAddress string `env:"RECALL_GRPC_SERVER_ADDRESS,required"`
}
