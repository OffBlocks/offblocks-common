package types

import (
	common "buf.build/gen/go/offblocks/offblocks-proto/protocolbuffers/go/common/v1"
	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func (m UUID) MarshalProto() (*common.UUID, error) {
	bytes, err := m.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &common.UUID{
		Uuid: bytes,
	}, nil
}

func (m *UUID) UnmarshalProto(pb *common.UUID) error {
	if pb == nil {
		return nil
	}
	var uuid UUID
	err := uuid.UnmarshalBinary(pb.Uuid)
	if err != nil {
		return err
	}
	*m = uuid
	return nil
}
