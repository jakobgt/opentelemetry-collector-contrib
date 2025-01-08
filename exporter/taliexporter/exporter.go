package taliexporter

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

type traceExporter struct {
}

func (t *traceExporter) ConsumeTracesFunc(ctx context.Context, td ptrace.Traces) error {
	batchpersignal.SplitTraces(td)
	td.ResourceSpans().At(0).Resource().Attributes()
	// TODO: Initialize Tali client and make a builder pattern for
	// creating the bloom filter.
	// For every trace we should just upload this batch
	return nil
}

// func (t *traceExporter) Start(ctx context.Context, host component.Host) error {
// 	panic("foobar")
// 	return nil
// }

// func (t *traceExporter) Stop(ctx context.Context) error {
// 	return nil
// }
