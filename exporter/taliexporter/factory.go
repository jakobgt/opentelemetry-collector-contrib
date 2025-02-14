// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package taliexporter

import (
	"context"
	"fmt"

	"github.com/jakobgt/go-tali"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/internal/metadata"
)

//go:generate mdatagen metadata.yaml

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, metadata.TracesStability),
	)
}

func createTracesExporter(_ context.Context, params exporter.Settings, comCfg component.Config) (exporter.Traces, error) {
	cfg, ok := comCfg.(Config)
	if !ok {
		return nil, fmt.Errorf("config was not a taliexporter.Config: %#v", comCfg)
	}

	logger := params.TelemetrySettings.Logger

	var arg tali.OtelClientArgs

	if cfg.S3DevNullMode {
		logger.Warn("TaliExporter is starting in /dev/null mode which means all data is sent to /dev/null. Use only for testing")
		arg = tali.WithObjectStorageClient(tali.NewDevNullClient())
	} else {
		arg = tali.WithHeadlessMode(
			tali.HeadlessModeConfig{
				S3AccessKey: cfg.S3AccessKey,
				S3SecretKey: cfg.S3SecretKey,
				S3Endpoint:  cfg.S3Endpoint,
				S3UseSSL:    cfg.S3UseSSL,
			})
	}

	// TODO: Should the initialization of the tali client happen in start?
	client, err := tali.NewOtelTraceClient(arg)
	if err != nil {
		return nil, err
	}
	te := newExporter(client, logger)
	return exporterhelper.NewTraces(
		context.TODO(),
		params,
		cfg,
		te.ConsumeTracesFunc,
		//		exporterhelper.WithStart(te.start),
		exporterhelper.WithShutdown(te.stop),
		// TODO: Consider to use start and stop functions?
		// exporterhelper.WithStart(te.Start),
		// exporterhelper.WithShutdown(te.Stop),

		// TODO: Once we can, we should use the batcher option
		// exporterhelper.WithBatcher()
	)
}
