package tali

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	// Shutdown shuts down the Tali client.
	Shutdown(ctx context.Context) error
}

var _ Client = (*client)(nil)

type client struct {
	s3client    minio.Client
	s3Transport *http.Transport
	ztdEncoder  *zstd.Encoder
}

// NewClient returns a Tali client
// Make with options
func NewClient() (Client, error) {
	compressor, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic("could not create tali client. Error is " + err.Error())
	}
	endpoint := "localhost:9000"
	// Note local Minio secrets.
	accessKeyID := "ZfhAK4x5LvEgsF1IIVE9"
	secretAccessKey := "UQut09ZtU0NyE5wY6Mychbg0e2oJTin7qNJpDy9P"
	useSSL := false
	// We need to be able to close any connections at the transport layer, so manually creating the transport here

	tr, err := minio.DefaultTransport(useSSL)
	if err != nil {
		return nil, err
	}
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:    useSSL,
		Transport: tr,
	})
	if err != nil {
		return nil, err
	}

	return &client{
		s3client:    *minioClient,
		s3Transport: tr,
		ztdEncoder:  compressor,
	}, nil
}

// Shutdown closes all idle outgoing connections from the Tali client
func (c *client) Shutdown(ctx context.Context) error {
	c.s3Transport.CloseIdleConnections()
	return nil
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
	s3key = fmt.Sprintf("%s.log", s3key)
	contentType := "application/octet-stream"
	_, err = c.s3client.PutObject(ctx, S3_BUCKET, s3key, bytes.NewReader(payload), int64(len(payload)),
		minio.PutObjectOptions{ContentType: contentType})
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
