package types

import (
	"database/sql/driver"
	"fmt"
	"net/url"

	common "buf.build/gen/go/offblocks/offblocks-proto/protocolbuffers/go/common/v1"
	"github.com/offblocks/offblocks-common/util"
)

type URL struct {
	url.URL
}

func Parse(s string) (URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return URL{}, err
	}
	return URL{*u}, nil
}

func MustParse(s string) URL {
	u, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

// MarshalText implements the encoding.TextMarshaler interface for XML
func (m URL) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
func (m *URL) UnmarshalText(data []byte) error {
	url, err := Parse(string(data))
	if err != nil {
		return err
	}
	*m = url
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (m URL) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *URL) UnmarshalJSON(data []byte) error {
	str, err := util.UnquoteIfQuoted(data)
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", data, err)
	}

	url, err := Parse(str)
	if err != nil {
		return err
	}
	*m = url
	return nil
}

func (m URL) MarshalProto() (*common.URL, error) {
	return &common.URL{
		Url: m.String(),
	}, nil
}

func (m *URL) UnmarshalProto(pb *common.URL) error {
	if pb == nil {
		return nil
	}
	url, err := Parse(pb.Url)
	if err != nil {
		return err
	}
	*m = url
	return nil
}

func (u URL) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *URL) Scan(value interface{}) error {
	var us string
	switch v := value.(type) {
	case []byte:
		us = string(v)
	case string:
		us = v
	default:
		return fmt.Errorf("failed to parse URL: %s", value)
	}
	uu, err := Parse(us)
	if err != nil {
		return err
	}
	*u = uu
	return nil
}
