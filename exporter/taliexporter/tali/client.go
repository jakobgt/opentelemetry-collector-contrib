package tali

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali/ngram"
)

// Segment represents a tali segment that can be uploaded to a Tali destination.
type Segment struct {
	// Segment is really used to shield the underlying protobuf definition.
}

type Client interface {
	// Upload uploads the Tali segment to the right destination
	Upload(ctx context.Context, seg Segment) error

	GenerateSegment(ctx context.Context, builder ngram.Builder, bytz []byte) (Segment, error)
}

var _ Client = (*client)(nil)

type client struct {
	// TODO: Bloomfilter
}

// NewClient returns a Tali client
// Make with options
func NewClient() Client {
	return &client{}
}

// Upload implements Client.
func (c *client) Upload(ctx context.Context, seg Segment) error {
	return nil
}

// GenerateSegment implements Client.
func (c *client) GenerateSegment(ctx context.Context, builder ngram.Builder, bytz []byte) (Segment, error) {
	panic("unimplemented")
}
