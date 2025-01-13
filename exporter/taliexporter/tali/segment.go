package tali

import (
	talipb "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/taliexporter/tali/internal/gen_proto"
	"google.golang.org/protobuf/proto"
)

// Segment represents a tali segment that can be uploaded to a Tali destination.
type Segment struct {
	// Segment is really used to shield the underlying protobuf definition.
	pbsegment *talipb.Segment
}

// Marshal returns the underlying protobuf marshalled
func (s *Segment) Marshal() ([]byte, error) {
	return proto.Marshal(s.pbsegment)
}
