package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/offblocks/offblocks-common/types"
	"github.com/stretchr/testify/require"
)

func TestUUIDMarshalProto(t *testing.T) {
	for _, tc := range []struct {
		uuid uuid.UUID
	}{{
		uuid: uuid.New(),
	}, {
		uuid: uuid.NewMD5(uuid.NameSpaceOID, []byte("test")),
	}, {
		uuid: uuid.NewSHA1(uuid.NameSpaceOID, []byte("test")),
	}} {
		uuid := types.UUID{UUID: tc.uuid}
		pb, err := uuid.MarshalProto()
		if err != nil {
			t.Fatalf("Failed to marshal UUID: %v", err)
		}

		var unmarshaled types.UUID
		if err := unmarshaled.UnmarshalProto(pb); err != nil {
			t.Fatalf("Failed to unmarshal UUID: %v", err)
		}

		require.Equal(t, uuid, unmarshaled)
	}
}
