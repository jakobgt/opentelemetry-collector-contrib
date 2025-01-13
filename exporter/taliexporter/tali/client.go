package tali

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/klauspost/compress/zstd"
	talipb "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali/internal/gen_proto"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali/ngram"
)

const (
	S3_BUCKET = "tali"
)

type Client interface {
	// Upload uploads the Tali segment to the right destination
	Upload(ctx context.Context, seg Segment) error

	GenerateSegment(ctx context.Context, builder ngram.Builder, bytz []byte) (Segment, error)
}

var _ Client = (*client)(nil)

type client struct {
	s3Manager  *manager.Uploader
	ztdEncoder *zstd.Encoder
}

// NewClient returns a Tali client
// Make with options
func NewClient() (Client, error) {
	compressor, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic("could not create tali client. Error is " + err.Error())
	}

	// TODO: Should we allow configuration of the part size? 5MiB is small if running in the same
	// cloud zone as S3.
	cfg, err := config.LoadDefaultConfig(context.TODO()) // Replace with your AWS region
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
	})

	return &client{
		ztdEncoder: compressor,
		s3Manager:  uploader,
	}, nil
}

// Upload implements Client.
func (c *client) Upload(ctx context.Context, seg Segment) error {
	// TODO: Should we read the first byte to determine the version?
	payload, err := seg.Marshal()
	if err != nil {
		return err
	}
	// TODO: Add an actual call to the Tali meta server here.
	// TODO: This call should actually be with a signed URL.

	t := time.Now()
	s3key := t.Format("20060102150405")

	input := &s3.PutObjectInput{
		Bucket:            aws.String(S3_BUCKET),
		Key:               aws.String(fmt.Sprintf("/segment/%s.seg", s3key)),
		Body:              bytes.NewReader(payload),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
	}
	// TODO: If the segment is larger than 5Mib, then it becomes a Multipart.
	_, err = c.s3Manager.Upload(ctx, input)
	return err
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
