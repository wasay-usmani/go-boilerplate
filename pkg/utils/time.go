package utils

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

func NowInLocation(location string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Karachi")
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(loc), nil
}
