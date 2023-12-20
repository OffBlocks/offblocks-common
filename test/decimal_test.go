package test

import (
	"testing"

	"github.com/offblocks/offblocks-common/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestDecimalMarshalProto(t *testing.T) {
	for _, tc := range []struct {
		d decimal.Decimal
	}{{
		d: decimal.RequireFromString("-123.45"),
	}, {
		d: decimal.RequireFromString(".0001"),
	}, {
		d: decimal.RequireFromString("1.47000"),
	}} {
		d := types.Decimal{Decimal: tc.d}
		pb, err := d.MarshalProto()
		if err != nil {
			t.Fatalf("Failed to marshal decimal: %v", err)
		}

		var unmarshaled types.Decimal
		if err := unmarshaled.UnmarshalProto(pb); err != nil {
			t.Fatalf("Failed to unmarshal decimal: %v", err)
		}

		require.True(t, d.Equal(unmarshaled.Decimal))
	}
}
