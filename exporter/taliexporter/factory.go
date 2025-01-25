// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package taliexporter

import (
	"context"

	"github.com/jakobgt/go-tali"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/internal/metadata"

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

func createTracesExporter(_ context.Context, params exporter.Settings, cfg component.Config) (exporter.Traces, error) {
	// TODO: Should the initialization of the tali client happen in start?
	client, err := tali.NewOtelClient()
	if err != nil {
		return nil, err
	}
	te := newExporter(client)
	return exporterhelper.NewTraces(
		context.TODO(),
		params,
		cfg,
		te.ConsumeTracesFunc,
		exporterhelper.WithShutdown(client.Shutdown),
		// TODO: Consider to use start and stop functions?
		// exporterhelper.WithStart(te.Start),
		// exporterhelper.WithShutdown(te.Stop),
		// TODO: Once we can, we should use the batcher option
		//exporterhelper.WithBatcher()
	)
}
