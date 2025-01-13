package tali

import (
	"context"

	"github.com/klauspost/compress/zstd"
	talipb "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali/internal/gen_proto"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali/ngram"
)

// Segment represents a tali segment that can be uploaded to a Tali destination.
type Segment struct {
	// Segment is really used to shield the underlying protobuf definition.
	pbsegment *talipb.Segment
}

type Client interface {
	// Upload uploads the Tali segment to the right destination
	Upload(ctx context.Context, seg Segment) error

	GenerateSegment(ctx context.Context, builder ngram.Builder, bytz []byte) (Segment, error)
}

var _ Client = (*client)(nil)

type client struct {
	ztdEncoder *zstd.Encoder
	// TODO: Bloomfilter
}

// NewClient returns a Tali client
// Make with options
func NewClient() Client {
	compressor, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic("could not create tali client. Error is " + err.Error())
	}
	return &client{
		ztdEncoder: compressor,
	}
}

// Upload implements Client.
func (c *client) Upload(ctx context.Context, seg Segment) error {
	return nil
}

// GenerateSegment implements Client.
func (c *client) GenerateSegment(ctx context.Context, builder ngram.Builder, bytz []byte) (Segment, error) {
	ngramIndex, err := builder.Build()
	if err != nil {
		return Segment{}, nil
	}
	ngramBytz, err := ngramIndex.Marshal()
	if err != nil {
		return Segment{}, err
	}

	dst := []byte{}
	dst = c.ztdEncoder.EncodeAll(bytz, dst)

	inner := talipb.Segment{
		Config: &talipb.Segment_Config{
			CompressionMethod: talipb.Segment_ZSTD,
			NgramIndex:        talipb.Segment_BINARY_FUSE,
		},
		NgramIndex:     ngramBytz,
		CompressedData: dst,
	}

	return Segment{
		&inner,
	}, nil
}
