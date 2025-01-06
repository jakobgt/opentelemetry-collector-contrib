package taliexporter

import (
	"context"
	"open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/internal/metadata"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

//go:generate mdatagen metadata.yaml

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, metadata.TracesStability),
	)
}

type Config struct {
}

func createDefaultConfig() component.Config {
	return Config{}
}

func createTracesExporter(_ context.Context, params exporter.Settings, cfg component.Config) (exporter.Traces, error) {
	//exporterConfig := cfg.(*Config)
	return nil, nil
}
