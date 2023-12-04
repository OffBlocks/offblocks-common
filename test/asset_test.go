package test

import (
	"encoding/json"
	"testing"

	"github.com/offblocks/offblocks-common/blockchain"
)

// See: https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-19.md#test-cases
func TestAssetId(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// Ether Token
		id: "eip155:1/slip44:60",
	}, {
		// Bitcoin Token
		id: "bip122:000000000019d6689c085ae165831e93/slip44:0",
	}, {
		// ATOM Token
		id: "cosmos:cosmoshub-3/slip44:118",
	}, {
		// Litecoin Token
		id: "bip122:12a765e31ffd4059bada1e25190f6e98/slip44:2",
	}, {
		// Binance Token
		id: "cosmos:Binance-Chain-Tigris/slip44:714",
	}, {
		// IOV Token
		id: "cosmos:iov-mainnet/slip44:234",
	}, {
		// Lisk Token
		id: "lip9:9ee11e9df416b18b/slip44:134",
	}, {
		// DAI Token
		id: "eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0f",
	}, {
		// USDC (Solana) Token
		id: "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp/spl:EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}, {
		// CryptoKitties Collectible
		id: "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d",
	}, {
		// CryptoKitties Collectible ID
		id: "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d/771769",
	}} {
		a := blockchain.AssetId{}
		if err := a.Parse(tc.id); err != nil {
			t.Fatalf("Failed to parse asset id: %v", err)
		}

		if a.String() != tc.id {
			t.Fatalf("Failed to serialize asset id to string")
		}

		if _, err := blockchain.NewAssetId(a.ChainId, a.Namespace, a.Reference); err != nil {
			t.Fatalf("Failed to create asset id from namespace and reference")
		}

		b, err := a.MarshalText()
		if err != nil {
			t.Errorf("Failed to marshal to text")
		}

		a = blockchain.AssetId{}
		if err := a.UnmarshalText(b); err != nil {
			t.Errorf("Failed to unmarshal from text")
		}

		if a.String() != tc.id {
			t.Errorf("Unmarshalled asset id invalid")
		}

		b, err = json.Marshal(a)
		if err != nil {
			t.Fatalf("Failed to marshal to json")
		}

		a = blockchain.AssetId{}
		if err := json.Unmarshal(b, &a); err != nil {
			t.Fatalf("Failed to unmarshal to json")
		}

		if a.String() != tc.id {
			t.Fatalf("Unmarshalled asset id invalid")
		}

		a2 := blockchain.AssetId{}
		if err := a2.Scan(a.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if a2.String() != a.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}
