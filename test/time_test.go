package test

import (
	"testing"
	"time"

	"github.com/offblocks/offblocks-common/types"
	"github.com/stretchr/testify/require"
)

func TestTimeMarshalProto(t *testing.T) {
	for _, tc := range []struct {
		time time.Time
	}{{
		time: time.Now().UTC(),
	}} {
		time := types.Time{Time: tc.time}
		pb, err := time.MarshalProto()
		if err != nil {
			t.Fatalf("Failed to marshal decimal: %v", err)
		}

		var unmarshaled types.Time
		if err := unmarshaled.UnmarshalProto(pb); err != nil {
			t.Fatalf("Failed to unmarshal decimal: %v", err)
		}

		require.Equal(t, time, unmarshaled)
	}
}
