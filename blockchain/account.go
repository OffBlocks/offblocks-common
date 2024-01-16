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
	if err := aID.validate(); err != nil {
		return AccountId{}, err
	}

	return aID, nil
}

func (a AccountId) validate() error {
	if err := a.ChainId.validate(); err != nil {
		return err
	}

	if ok := addressRegex.Match([]byte(a.Address)); !ok {
		return errors.New("namespace does not match spec")
	}

	return nil
}

// String returns the string form of account id, chain_namespace:chain_reference:address
func (a AccountId) String() string {
	return a.ChainId.String() + ":" + a.Address
}

// Parse parses a string into a account id from the string form, chain_namespace:chain_reference:address
func (a *AccountId) Parse(s string) error {
	split := strings.SplitN(s, ":", 3)
	if len(split) != 3 {
		return fmt.Errorf("invalid account id: %s", s)
	}

	*a = AccountId{ChainId{split[0], split[1]}, split[2]}
	if err := a.validate(); err != nil {
		return err
	}

	return nil
}

// MustParse parses a string into a account id from the string form, chain_namespace:chain_reference:address
// and panics if there is an error
func (c *AccountId) MustParse(s string) {
	if err := c.Parse(s); err != nil {
		panic(err)
	}
}

// ParseAccountId parses a string into a account id from the string form, chain_namespace:chain_reference:address
func ParseAccountId(s string) (AccountId, error) {
	var a AccountId
	err := a.Parse(s)
	if err != nil {
		return a, err
	}

	return a, nil
}

// MustParseAccountId parses a string into a account id from the string form, chain_namespace:chain_reference:address
// and panics if there is an error
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
