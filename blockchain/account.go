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

type AccountId struct {
	ChainId ChainId
	Address string
}

var (
	addressRegex = regexp.MustCompile("[a-zA-Z0-9]{1,64}")
)

func NewAccountId(chainId ChainId, address string) (AccountId, error) {
	aID := AccountId{chainId, address}
	if err := aID.Validate(); err != nil {
		return AccountId{}, err
	}

	return aID, nil
}

func UnsafeAccountId(chainId ChainId, address string) AccountId {
	return AccountId{chainId, address}
}

func (a AccountId) Validate() error {
	if err := a.ChainId.Validate(); err != nil {
		return err
	}

	if ok := addressRegex.Match([]byte(a.Address)); !ok {
		return errors.New("namespace does not match spec")
	}

	return nil
}

func (a AccountId) String() string {
	if err := a.Validate(); err != nil {
		panic(err)
	}
	return a.ChainId.String() + ":" + a.Address
}

func (a *AccountId) Parse(s string) error {
	split := strings.SplitN(s, ":", 3)
	if len(split) != 3 {
		return fmt.Errorf("invalid account id: %s", s)
	}

	*a = AccountId{ChainId{split[0], split[1]}, split[2]}
	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *AccountId) MustParse(s string) {
	if err := c.Parse(s); err != nil {
		panic(err)
	}
}

func ParseAccountId(s string) (AccountId, error) {
	var a AccountId
	err := a.Parse(s)
	if err != nil {
		return a, err
	}

	return a, nil
}

func MustParseAccountId(s string) AccountId {
	var a AccountId
	a.MustParse(s)
	return a
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization
func (a *AccountId) UnmarshalText(data []byte) error {
	accountId, err := ParseAccountId(string(data))
	if err != nil {
		return err
	}
	*a = accountId
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization
func (a AccountId) MarshalText() ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	return []byte(a.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *AccountId) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	str, err := util.UnquoteIfQuoted(data)
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", data, err)
	}

	accountId, err := ParseAccountId(str)
	if err != nil {
		return err
	}
	*a = accountId
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (a AccountId) MarshalJSON() ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	str := "\"" + a.String() + "\""

	return []byte(str), nil
}

func (a *AccountId) UnmarshalProto(pb string) error {
	accountId, err := ParseAccountId(pb)
	if err != nil {
		return err
	}
	*a = accountId
	return nil
}

func (a AccountId) MarshalProto() (string, error) {
	return a.String(), nil
}

func (a AccountId) Value() (driver.Value, error) {
	return a.String(), nil
}

func (a *AccountId) Scan(src interface{}) error {
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

func (a AccountId) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(strings.ToUpper(a.String())))
}

func (a *AccountId) UnmarshalGQL(v interface{}) error {
	if id, ok := v.(string); ok {
		if err := a.Parse(id); err != nil {
			return fmt.Errorf("unmarshalling account id: %w", err)
		}
	}

	return nil
}
