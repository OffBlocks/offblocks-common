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

type TransactionId struct {
	ChainId ChainId
	Hash    string
}

var (
	hashRegex = regexp.MustCompile("[a-zA-Z0-9]{1,128}")
)

func NewTransactionId(ChainId ChainId, hash string) (TransactionId, error) {
	aID := TransactionId{ChainId, hash}
	if err := aID.Validate(); err != nil {
		return TransactionId{}, err
	}

	return aID, nil
}

func UnsafeTransactionId(chainId ChainId, hash string) TransactionId {
	return TransactionId{chainId, hash}
}

func (t TransactionId) Validate() error {
	if err := t.ChainId.Validate(); err != nil {
		return err
	}

	if ok := hashRegex.Match([]byte(t.Hash)); !ok {
		return errors.New("namespace does not match spec")
	}

	return nil
}

func (t TransactionId) String() string {
	if err := t.Validate(); err != nil {
		panic(err)
	}
	return t.ChainId.String() + ":" + t.Hash
}

func (t *TransactionId) Parse(s string) error {
	split := strings.SplitN(s, ":", 3)
	if len(split) != 3 {
		return fmt.Errorf("invalid transaction id: %s", s)
	}

	*t = TransactionId{ChainId{split[0], split[1]}, split[2]}
	if err := t.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *TransactionId) MustParse(s string) {
	if err := c.Parse(s); err != nil {
		panic(err)
	}
}

func ParseTransactionId(s string) (TransactionId, error) {
	var t TransactionId
	err := t.Parse(s)
	if err != nil {
		return t, err
	}

	return t, nil
}

func MustParseTransactionId(s string) TransactionId {
	var t TransactionId
	t.MustParse(s)
	return t
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization
func (t *TransactionId) UnmarshalText(data []byte) error {
	TransactionId, err := ParseTransactionId(string(data))
	if err != nil {
		return err
	}
	*t = TransactionId
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization
func (t TransactionId) MarshalText() ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	return []byte(t.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionId) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	str, err := unquoteIfQuoted(data)
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", data, err)
	}

	transactionId, err := ParseTransactionId(str)
	if err != nil {
		return err
	}
	*t = transactionId
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t TransactionId) MarshalJSON() ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	str := "\"" + t.String() + "\""

	return []byte(str), nil
}

func (t TransactionId) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *TransactionId) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return fmt.Errorf("scanning transaction id: %w", err)
	}

	if !i.Valid {
		return nil
	}

	if err := t.Parse(i.String); err != nil {
		return err
	}

	return nil
}

func (t TransactionId) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(strings.ToUpper(t.String())))
}

func (t *TransactionId) UnmarshalGQL(v interface{}) error {
	if id, ok := v.(string); ok {
		if err := t.Parse(id); err != nil {
			return fmt.Errorf("unmarshalling transaction id: %w", err)
		}
	}

	return nil
}
