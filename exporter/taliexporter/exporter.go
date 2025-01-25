// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package taliexporter

import (
	"context"

	"github.com/jakobgt/go-tali"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type traceExporter struct {
	tclient        tali.Client[ptrace.Traces]
	traceMarshaler ptrace.Marshaler
	// TODO: Add for logs (and metrics?) as well
}

func newExporter(client tali.Client[ptrace.Traces]) *traceExporter {
	return &traceExporter{
		tclient:        client,
		traceMarshaler: &ptrace.JSONMarshaler{},
	}
}

func (t *traceExporter) ConsumeTracesFunc(ctx context.Context, td ptrace.Traces) error {
	// batchpersignal.SplitTraces(td)
	// TODO: Batch traces per traceid
	return t.tclient.Upload(ctx, td)
}
