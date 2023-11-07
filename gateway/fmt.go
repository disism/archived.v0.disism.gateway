package gateway

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// DateStamp returns a timestamppb.Timestamp object based on the input time.Time.
//
// It takes a time.Time parameter and returns a *timestamppb.Timestamp object.
func DateStamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
