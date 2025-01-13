package taliexporter

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
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
	// TODO: Create the tali client.
	client := tali.NewClient()
	te := newExporter(client)
	return exporterhelper.NewTraces(
		context.TODO(),
		params,
		cfg,
		te.ConsumeTracesFunc,
		// TODO: Consider to use start and stop functions?
		// exporterhelper.WithStart(te.Start),
		// exporterhelper.WithShutdown(te.Stop),
		// TODO: Once we can, we should use the batcher option
		//exporterhelper.WithBatcher()
	)
}
