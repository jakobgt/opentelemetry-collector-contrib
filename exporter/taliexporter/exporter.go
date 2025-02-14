// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package taliexporter

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/jakobgt/go-tali"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type traceExporter struct {
	traceClient tali.Client[ptrace.Traces]

	stopped    atomic.Bool
	goroutines sync.WaitGroup

	log *zap.Logger
	// TODO: Add for logs (and metrics?) as well
}

func newExporter(client tali.Client[ptrace.Traces], logger *zap.Logger) *traceExporter {
	return &traceExporter{
		traceClient: client,
		log:         logger,
		goroutines:  sync.WaitGroup{},
	}
}

func (t *traceExporter) ConsumeTracesFunc(ctx context.Context, td ptrace.Traces) error {
	// TODO: This is a blocking call from the batcher...

	// Consider whether we should split the incoming traces.

	// TODO: Consider to use the batching operator in the exporthelper config

	// Should we have a set of goroutines that just handle the upload part?
	// In some sense, we should just generate the segment in this one goroutine and then
	// do the upload in a background goroutine?

	//batchpersignal.SplitTraces(td)
	// TODO: Batch traces per traceid - Use the sort function on resourceSpans?

	// If we're stopped we just upload this one directly in this goroutine
	if t.stopped.Load() {
		return t.traceClient.Upload(ctx, td)
	}

	// otherwise we stop in a different one
	t.goroutines.Add(1)
	go func() {
		defer t.goroutines.Done()
		spanCount := td.SpanCount()
		t.log.Info("Uploading new Tali segment", zap.Int("span_count", spanCount))
		ctx := context.Background()
		if err := t.traceClient.Upload(ctx, td); err != nil {
			t.log.Warn("error encountered when uploading segment", zap.Int("span_count", spanCount), zap.Error(err))
		}
	}()
	return nil
}

// func (t *traceExporter) start(context.Context, component.Host) error {
// 	for
// }

func (t *traceExporter) stop(ctx context.Context) error {
	err := t.traceClient.Shutdown(ctx)
	t.stopped.Store(true)
	t.goroutines.Wait()
	return err
}
