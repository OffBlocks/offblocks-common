package blockchain

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type AssetId struct {
	ChainId   ChainId
	Namespace string
	Reference string
}

var (
	assetNamespaceRegex = regexp.MustCompile("[-a-z0-9]{3,8}")
	assetReferenceRegex = regexp.MustCompile("[-a-zA-Z0-9]{1,64}")
)

func NewAssetId(chainID ChainId, namespace, reference string) (AssetId, error) {
	aID := AssetId{chainID, namespace, reference}
	if err := aID.Validate(); err != nil {
		return AssetId{}, err
	}

	return aID, nil
}

func UnsafeAssetId(chainID ChainId, namespace, reference string) AssetId {
	return AssetId{chainID, namespace, reference}
}

func (a AssetId) Validate() error {
	if ok := assetNamespaceRegex.Match([]byte(a.Namespace)); !ok {
		return errors.New("asset namespace does not match spec")
	}

	if ok := assetReferenceRegex.Match([]byte(a.Reference)); !ok {
		return errors.New("asset reference does not match spec")
	}

	return nil
}

func (a AssetId) String() string {
	if err := a.Validate(); err != nil {
		panic(err)
	}
	return a.ChainId.String() + "/" + a.Namespace + ":" + a.Reference
}

func (a *AssetId) Parse(s string) error {
	components := strings.SplitN(s, "/", 2)
	if len(components) != 2 {
		return fmt.Errorf("invalid asset id: %s", s)
	}

	cID := new(ChainId)
	if err := cID.Parse(components[0]); err != nil {
		return err
	}

	asset := strings.SplitN(components[1], ":", 2)
	if len(asset) != 2 {
		return fmt.Errorf("invalid asset id: %s", s)
	}

	*a = AssetId{*cID, asset[0], asset[1]}
	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (a *AssetId) MustParse(s string) {
	if err := a.Parse(s); err != nil {
		panic(err)
	}
}

func ParseAssetId(s string) (AssetId, error) {
	var a AssetId
	err := a.Parse(s)
	if err != nil {
		return a, err
	}

	return a, nil
}

func MustParseAssetId(s string) AssetId {
	var a AssetId
	a.MustParse(s)
	return a
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization
func (a *AssetId) UnmarshalText(data []byte) error {
	assetId, err := ParseAssetId(string(data))
	if err != nil {
		return err
	}
	*a = assetId
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization
func (a AssetId) MarshalText() ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	return []byte(a.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *AssetId) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	str, err := unquoteIfQuoted(data)
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", data, err)
	}

	assetId, err := ParseAssetId(str)
	if err != nil {
		return err
	}
	*a = assetId
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (a AssetId) MarshalJSON() ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	str := "\"" + a.String() + "\""

	return []byte(str), nil
}

func (a AssetId) Value() (driver.Value, error) {
	return a.String(), nil
}

func (a *AssetId) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return fmt.Errorf("scanning account id: %w", err)
	}

	if !i.Valid {
		return nil
	}

	if err := a.Parse(i.String); err != nil {
		return err
	}

	return nil
}

func (a AssetId) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(strings.ToUpper(a.String())))
}

func (a *AssetId) UnmarshalGQL(v interface{}) error {
	if id, ok := v.(string); ok {
		if err := a.Parse(id); err != nil {
			return fmt.Errorf("unmarshalling asset id: %w", err)
		}
	}

	return nil
}
