package blockchain

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ChainId struct {
	Namespace string `json:"namespace"`
	Reference string `json:"reference"`
}

var (
	chainNamespaceRegex = regexp.MustCompile("[-a-z0-9]{3,8}")
	chainReferenceRegex = regexp.MustCompile("[-a-zA-Z0-9]{1,32}")
)

func NewChainId(namespace, reference string) (ChainId, error) {
	cID := ChainId{namespace, reference}
	if err := cID.Validate(); err != nil {
		return ChainId{}, err
	}

	return cID, nil
}

func UnsafeChainId(namespace, reference string) ChainId {
	return ChainId{namespace, reference}
}

func (c ChainId) Validate() error {
	if ok := chainNamespaceRegex.Match([]byte(c.Namespace)); !ok {
		return errors.New("chain namespace does not match spec")
	}

	if ok := chainReferenceRegex.Match([]byte(c.Reference)); !ok {
		return errors.New("chain reference does not match spec")
	}

	return nil
}

func (c ChainId) String() string {
	if err := c.Validate(); err != nil {
		panic(err)
	}
	return c.Namespace + ":" + c.Reference
}

func (c *ChainId) Parse(s string) error {
	split := strings.SplitN(s, ":", 2)
	if len(split) != 2 {
		return fmt.Errorf("invalid chain id: %s", s)
	}

	*c = ChainId{split[0], split[1]}
	if err := c.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *ChainId) MustParse(s string) {
	if err := c.Parse(s); err != nil {
		panic(err)
	}
}

func ParseChainId(s string) (ChainId, error) {
	var c ChainId
	err := c.Parse(s)
	if err != nil {
		return c, err
	}

	return c, nil
}

func MustParseChainId(s string) ChainId {
	var c ChainId
	c.MustParse(s)
	return c
}

func (c ChainId) MarshalText() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return []byte(c.String()), nil
}

func (c *ChainId) UnmarshalText(data []byte) error {
	chainId, err := ParseChainId(string(data))
	if err != nil {
		return err
	}
	*c = chainId
	return nil
}

func (c ChainId) MarshalJSON() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	type ChainIdAlias ChainId
	ca := (ChainIdAlias)(c)
	return json.Marshal(ca)
}

func (c *ChainId) UnmarshalJSON(data []byte) error {
	type ChainIdAlias ChainId
	ca := (*ChainIdAlias)(c)
	if err := json.Unmarshal(data, &ca); err != nil {
		return err
	}

	if err := c.Validate(); err != nil {
		return err
	}

	return nil
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
