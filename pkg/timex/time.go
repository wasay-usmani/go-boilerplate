package timex

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Now() time.Time {
	return time.Now().UTC()
}

func TimestampNow() *timestamppb.Timestamp {
	return timestamppb.New(Now())
}

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
