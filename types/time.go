package types

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Time struct {
	time.Time
}

func (m *Time) MarshalProto() (*timestamppb.Timestamp, error) {
	if m == nil {
		return nil, nil
	}
	pb := timestamppb.New(m.Time.UTC())
	return pb, nil
}

func (m *Time) UnmarshalProto(pb *timestamppb.Timestamp) error {
	if pb == nil {
		return nil
	}
	time := pb.AsTime().UTC()
	*m = Time{time}
	return nil
}
