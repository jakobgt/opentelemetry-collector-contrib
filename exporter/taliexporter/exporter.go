package taliexporter

import (
	"context"
	"strconv"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali/ngram"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type traceExporter struct {
	tclient        tali.Client
	traceMarshaler ptrace.Marshaler
	// TODO: Add for logs (and metrics?) as well
}

func newExporter(client tali.Client) *traceExporter {
	return &traceExporter{
		tclient:        client,
		traceMarshaler: &ptrace.JSONMarshaler{},
	}
}

func (t *traceExporter) ConsumeTracesFunc(ctx context.Context, td ptrace.Traces) error {
	//	batchpersignal.SplitTraces(td)
	// Make ngram
	// Consider to make ngram generation while marshalling the JSON?
	// TODO: Gather start and endtime.
	b := generateTaliNGram(td)
	bytz, err := t.traceMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}

	s, err := t.tclient.GenerateSegment(ctx, b, bytz)
	if err != nil {
		return err
	}
	// Upload
	return t.tclient.Upload(ctx, s)
}

func generateTaliNGram(td ptrace.Traces) ngram.Builder {
	b := ngram.NewBuilder(1024)
	resourceSpans := td.ResourceSpans()
	if resourceSpans.Len() == 0 {
		return b
	}

	for i := 0; i < resourceSpans.Len(); i++ {
		rs := resourceSpans.At(i)
		addValues(&b, rs.Resource().Attributes())

		ilss := rs.ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			//library := ils.Scope()

			spans := ils.Spans()
			for k := 0; k < spans.Len(); k++ {
				otelSpan := spans.At(k)
				addValues(&b, otelSpan.Attributes())
			}
		}
	}
	return b
}

func addValues(b *ngram.Builder, attrs pcommon.Map) {

	attrs.Range(func(key string, attr pcommon.Value) bool {
		switch attr.Type() {
		case pcommon.ValueTypeStr:
			b.Add(key, attr.Str())
		case pcommon.ValueTypeBool:
			b.Add(key, strconv.FormatBool(attr.Bool()))
		case pcommon.ValueTypeDouble:
			b.Add(key, strconv.FormatFloat(attr.Double(), 'g', -1, 64))
		case pcommon.ValueTypeInt:
			b.Add(key, strconv.FormatInt(attr.Int(), 10))
		case pcommon.ValueTypeEmpty:
			// TODO: Add null values?
			b.Add(key, "NULL")
		case pcommon.ValueTypeMap:
			// TODO: Should we recurse here?
		case pcommon.ValueTypeSlice:
			// TODO: Should we recurse here?
		case pcommon.ValueTypeBytes:
		}
		return true
	})
}

// func (t *traceExporter) Start(ctx context.Context, host component.Host) error {
// 	panic("foobar")
// 	return nil
// }

// func (t *traceExporter) Stop(ctx context.Context) error {
// 	return nil
// }
