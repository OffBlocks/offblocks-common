package test

import (
	"encoding/json"
	"testing"

	"github.com/offblocks/offblocks-common/blockchain"
)

// See: https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-10.md#test-cases
func TestAccountId(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// Ethereum mainnet
		id: "eip155:1:0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb",
	}, {
		// Solana mainnet
		id: "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp:7S3P4HxJpyyigGzodYwHtCxZyUQe9JiBMHyRWXArAaKv",
	}, {
		// Bitcoin mainnet
		id: "bip122:000000000019d6689c085ae165831e93:128Lkh3S7CkDTBZ8W7BbpsN3YYizJMp8p6",
	}, {
		// Cosmos Hub
		id: "cosmos:cosmoshub-3:cosmos1t2uflqwqe0fsj0shcfkrvpukewcw40yjj6hdc0",
	}, {
		// Kusama network
		id: "polkadot:b0a8d493285c2df73290dfb7e61f870f:5hmuyxw9xdgbpptgypokw4thfyoe3ryenebr381z9iaegmfy",
	}, {
		// Dummy max length (64+1+8+1+32 = 106 chars/bytes)
		id: "chainstd:8c3444cf8970a9e41a706fab93e7a6c4:9IU9l4BzmRdU8V03BugERXt6che9H2Ntu6f12KHiym9V0dl4me3p9pQNhmUbNlru",
	}} {
		a := blockchain.AccountId{}
		if err := a.Parse(tc.id); err != nil {
			t.Errorf("Failed to parse account id")
		}

		if a.String() != tc.id {
			t.Errorf("Failed to serialize account id to string")
		}

		if _, err := blockchain.NewAccountId(a.ChainId, a.Address); err != nil {
			t.Errorf("Failed to create account id from address")
		}

		b, err := a.MarshalText()
		if err != nil {
			t.Errorf("Failed to marshal to text")
		}

		a = blockchain.AccountId{}
		if err := a.UnmarshalText(b); err != nil {
			t.Errorf("Failed to unmarshal from text")
		}

		if a.String() != tc.id {
			t.Errorf("Unmarshalled account id invalid")
		}

		b, err = json.Marshal(a)
		if err != nil {
			t.Errorf("Failed to marshal to json")
		}

		a = blockchain.AccountId{}
		if err := json.Unmarshal(b, &a); err != nil {
			t.Errorf("Failed to unmarshal to json")
		}

		pb, err := a.MarshalProto()
		if err != nil {
			t.Errorf("Failed to marshal to proto")
		}

		a = blockchain.AccountId{}
		if err := a.UnmarshalProto(pb); err != nil {
			t.Errorf("Failed to unmarshal from proto")
		}

		if a.String() != tc.id {
			t.Errorf("Unmarshalled account id invalid")
		}

		a2 := blockchain.AccountId{}
		if err := a2.Scan(a.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if a2.String() != a.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}
