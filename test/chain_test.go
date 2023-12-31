package test

import (
	"encoding/json"
	"testing"

	"github.com/offblocks/offblocks-common/blockchain"
)

// See: https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-2.md#test-cases
func TestChainId(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// Ethereum mainnet
		id: "eip155:1",
	}, {
		// Solana mainnet
		id: "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp",
	}, {
		// Bitcoin mainnet (see https://github.com/bitcoin/bips/blob/master/bip-0122.mediawiki#definition-of-chain-id)
		id: "bip122:000000000019d6689c085ae165831e93",
	}, {
		// Litecoin
		id: "bip122:12a765e31ffd4059bada1e25190f6e98",
	}, {
		// Feathercoin (Litecoin fork)
		id: "bip122:fdbe99b90c90bae7505796461471d89a",
	}, {
		// Cosmos Hub (Tendermint + Cosmos SDK)
		id: "cosmos:cosmoshub-2",
	}, {
		// Cosmos Hub (Tendermint + Cosmos SDK)
		id: "cosmos:cosmoshub-3",
	}, {
		// Binance chain (Tendermint + Cosmos SDK; see https://dataseed5.defibit.io/genesis)
		id: "cosmos:Binance-Chain-Tigris",
	}, {
		// IOV Mainnet (Tendermint + weave)
		id: "cosmos:iov-mainnet",
	}, {
		// Lisk Mainnet (LIP-0009; see https://github.com/LiskHQ/lips/blob/master/proposals/lip-0009.md)
		id: "lip9:9ee11e9df416b18b",
	}, {
		// Dummy max length (8+1+32 = 41 chars/bytes)
		id: "chainstd:8c3444cf8970a9e41a706fab93e7a6c4",
	}} {
		c := blockchain.ChainId{}
		if err := c.Parse(tc.id); err != nil {
			t.Errorf("Failed to parse chain id")
		}

		if c.String() != tc.id {
			t.Errorf("Failed to serialize chain id to string")
		}

		if _, err := blockchain.NewChainId(c.Namespace, c.Reference); err != nil {
			t.Errorf("Failed to create chain id from namespace and reference")
		}

		b, err := c.MarshalText()
		if err != nil {
			t.Errorf("Failed to marshal to text")
		}

		c = blockchain.ChainId{}
		if err := c.UnmarshalText(b); err != nil {
			t.Errorf("Failed to unmarshal from text")
		}

		if c.String() != tc.id {
			t.Errorf("Unmarshalled chain id invalid")
		}

		b, err = json.Marshal(c)
		if err != nil {
			t.Errorf("Failed to marshal to json")
		}

		c = blockchain.ChainId{}
		if err := json.Unmarshal(b, &c); err != nil {
			t.Errorf("Failed to unmarshal from json")
		}

		pb, err := c.MarshalProto()
		if err != nil {
			t.Errorf("Failed to marshal to proto")
		}

		c = blockchain.ChainId{}
		if err := c.UnmarshalProto(pb); err != nil {
			t.Errorf("Failed to unmarshal from proto")
		}

		if c.String() != tc.id {
			t.Errorf("Unmarshalled chain id invalid")
		}

		c2 := blockchain.ChainId{}
		if err := c2.Scan(c.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if c2.String() != c.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}
