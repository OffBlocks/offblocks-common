package test

import (
	"encoding/json"
	"testing"

	"github.com/offblocks/offblocks-common/blockchain"
)

// See: https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-10.md#test-cases
func TestTransactionId(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// Ethereum mainnet
		id: "eip155:1:0x66f2462a072d837b5c4a76de103a7e5d1cd42c5f77fbd4f95a0dcc9fddf90b08",
	}, {
		// Solana mainnet
		id: "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp:3QorZbZ5bRWePAqRsAW5vMLggKfaJQ7RTdKEoXQBsWtgunjrQAN8wrj99yLDqAQassU7DzVYdB62rygzKQd7m7fU",
	}, {
		// Bitcoin mainnet
		id: "bip122:000000000019d6689c085ae165831e93:c55e6d98f3867f5bffdd3fae24082ba56a50e81e13c46b67716343a1fedda9ba",
	}, {
		// Cosmos Hub
		id: "cosmos:cosmoshub-3:A57352B805703E81164196D050D9DBAC3283304518A421CD0BC4767C143E02ED",
	}, {
		// Kusama network
		id: "polkadot:b0a8d493285c2df73290dfb7e61f870f:0x87232efe499130a032cceed485e7bd54f22cfbd92477bd1b09f7d6d3dc1b7c1d",
	}, {
		// Dummy max length (128+1+8+1+32 = 170 chars/bytes)
		id: "chainstd:8c3444cf8970a9e41a706fab93e7a6c4:wwZRxE3frjZT988lmpUh6LIl31oXJTm1x2RU1ruu3RQkZ5wdvBT8FG0zOPFHGOEDAynr6SPWL7wXhHTipFQY6xBGTaPNjKpPNH4KTHCnEu0RqCffIS3ZRY5X85SNZV1O",
	}} {
		tt := blockchain.TransactionId{}
		if err := tt.Parse(tc.id); err != nil {
			t.Errorf("Failed to parse transaction id")
		}

		if tt.String() != tc.id {
			t.Errorf("Failed to serialize transaction id to string")
		}

		if _, err := blockchain.NewTransactionId(tt.ChainId, tt.Hash); err != nil {
			t.Errorf("Failed to create transaction id from address")
		}

		b, err := tt.MarshalText()
		if err != nil {
			t.Errorf("Failed to marshal to text")
		}

		tt = blockchain.TransactionId{}
		if err := tt.UnmarshalText(b); err != nil {
			t.Errorf("Failed to unmarshal from text")
		}

		if tt.String() != tc.id {
			t.Errorf("Unmarshalled transaction id invalid")
		}

		b, err = json.Marshal(tt)
		if err != nil {
			t.Errorf("Failed to marshal to json")
		}

		tt = blockchain.TransactionId{}
		if err := json.Unmarshal(b, &tt); err != nil {
			t.Errorf("Failed to unmarshal to json")
		}

		if tt.String() != tc.id {
			t.Errorf("Unmarshalled account id invalid")
		}

		a2 := blockchain.TransactionId{}
		if err := a2.Scan(tt.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if a2.String() != tt.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}
