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

	"github.com/offblocks/offblocks-common/util"
)

type ChainId struct {
	Namespace string
	Reference string
}

var (
	chainNamespaceRegex = regexp.MustCompile("[-a-z0-9]{3,8}")
	chainReferenceRegex = regexp.MustCompile("[-a-zA-Z0-9]{1,32}")
)

func NewChainId(namespace, reference string) (ChainId, error) {
	cID := ChainId{namespace, reference}
	if err := cID.validate(); err != nil {
		return ChainId{}, err
	}

	return cID, nil
}

func (c ChainId) validate() error {
	if ok := chainNamespaceRegex.Match([]byte(c.Namespace)); !ok {
		return errors.New("chain namespace does not match spec")
	}

	if ok := chainReferenceRegex.Match([]byte(c.Reference)); !ok {
		return errors.New("chain reference does not match spec")
	}

	return nil
}

// String returns the string form of chain id, namespace:reference
func (c ChainId) String() string {
	return c.Namespace + ":" + c.Reference
}

// Parse parses a string into a chain id from the string form, namespace:reference
func (c *ChainId) Parse(s string) error {
	split := strings.SplitN(s, ":", 2)
	if len(split) != 2 {
		return fmt.Errorf("invalid chain id: %s", s)
	}

	*c = ChainId{split[0], split[1]}
	if err := c.validate(); err != nil {
		return err
	}

	return nil
}

// MustParse parses a string into a chain id from the string form, namespace:reference
// and panics if there is an error
func (c *ChainId) MustParse(s string) {
	if err := c.Parse(s); err != nil {
		panic(err)
	}
}

// ParseChainId parses a string into a chain id from the string form, namespace:reference
func ParseChainId(s string) (ChainId, error) {
	var c ChainId
	err := c.Parse(s)
	if err != nil {
		return c, err
	}

	return c, nil
}

// MustParseChainId parses a string into a chain id from the string form, namespace:reference
// and panics if there is an error
func MustParseChainId(s string) ChainId {
	var c ChainId
	c.MustParse(s)
	return c
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization
func (c *ChainId) UnmarshalText(data []byte) error {
	chainId, err := ParseChainId(string(data))
	if err != nil {
		return err
	}
	*c = chainId
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization
func (c ChainId) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *ChainId) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	str, err := util.UnquoteIfQuoted(data)
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", data, err)
	}

	chainId, err := ParseChainId(str)
	if err != nil {
		return err
	}
	*c = chainId
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (c ChainId) MarshalJSON() ([]byte, error) {
	str := "\"" + c.String() + "\""

	return []byte(str), nil
}

func (c *ChainId) UnmarshalProto(pb string) error {
	chainId, err := ParseChainId(pb)
	if err != nil {
		return err
	}
	*c = chainId
	return nil
}

func (c ChainId) MarshalProto() (string, error) {
	return c.String(), nil
}

func (c ChainId) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *ChainId) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return fmt.Errorf("scanning chain id: %w", err)
	}

	if !i.Valid {
		return nil
	}

	if err := c.Parse(i.String); err != nil {
		return err
	}

	return nil
}

func (c ChainId) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(strings.ToUpper(c.String())))
}

func (c *ChainId) UnmarshalGQL(v interface{}) error {
	if id, ok := v.(string); ok {
		if err := c.Parse(id); err != nil {
			return fmt.Errorf("unmarshalling account id: %w", err)
		}
	}

	return nil
}
