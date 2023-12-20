package types

import (
	common "buf.build/gen/go/offblocks/offblocks-proto/protocolbuffers/go/common/v1"
	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (m Decimal) MarshalProto() (*common.Decimal, error) {
	return &common.Decimal{
		Decimal: m.String(),
	}, nil
}

func (m *Decimal) UnmarshalProto(pb *common.Decimal) error {
	if pb == nil {
		return nil
	}
	decimal, err := decimal.NewFromString(pb.Decimal)
	if err != nil {
		return err
	}
	*m = Decimal{decimal}
	return nil
}
