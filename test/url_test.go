package test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/offblocks/offblocks-common/types"
	"github.com/stretchr/testify/require"
)

func MustParse(s string) url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return *u
}

func TestURLMarshalProto(t *testing.T) {
	for _, tc := range []struct {
		url url.URL
	}{{
		url: MustParse("https://example.com/foo"),
	}, {
		url: MustParse("http://example.com/foo"),
	}, {
		url: MustParse("http://localhost:8080"),
	}, {
		url: MustParse("https://localhost:8080"),
	}} {
		url := types.URL{URL: tc.url}
		pb, err := url.MarshalProto()
		if err != nil {
			t.Fatalf("Failed to marshal URL: %v", err)
		}

		var unmarshaled types.URL
		if err := unmarshaled.UnmarshalProto(pb); err != nil {
			t.Fatalf("Failed to unmarshal URL: %v", err)
		}

		require.Equal(t, url, unmarshaled)
	}
}

func TestURLMarshalText(t *testing.T) {
	for _, tc := range []struct {
		url url.URL
	}{{
		url: MustParse("https://example.com/foo"),
	}, {
		url: MustParse("http://example.com/foo"),
	}, {
		url: MustParse("http://localhost:8080"),
	}, {
		url: MustParse("https://localhost:8080"),
	}} {
		url := types.URL{URL: tc.url}
		text, err := url.MarshalText()
		if err != nil {
			t.Fatalf("Failed to marshal URL: %v", err)
		}

		var unmarshaled types.URL
		if err := unmarshaled.UnmarshalText(text); err != nil {
			t.Fatalf("Failed to unmarshal URL: %v", err)
		}

		require.Equal(t, url, unmarshaled)
	}
}

func TestURLMarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		url url.URL
	}{{
		url: MustParse("https://example.com/foo"),
	}, {
		url: MustParse("http://example.com/foo"),
	}, {
		url: MustParse("http://localhost:8080"),
	}, {
		url: MustParse("https://localhost:8080"),
	}} {
		url := types.URL{URL: tc.url}
		j, err := json.Marshal(url)
		if err != nil {
			t.Fatalf("Failed to marshal URL: %v", err)
		}

		var unmarshaled types.URL
		if err := json.Unmarshal(j, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal URL: %v", err)
		}

		require.Equal(t, url, unmarshaled)
	}
}
