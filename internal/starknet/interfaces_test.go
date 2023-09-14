package starknet

import (
	"os"
	"testing"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestInterfaces(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantLen int
		wantErr bool
	}{
		{
			name:    "load interfaces",
			dir:     "../../build/interfaces",
			wantLen: 13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Interfaces(tt.dir)
			require.NoError(t, err)
			require.Len(t, got, tt.wantLen, "interfaces count")
		})
	}
}

func TestFindInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     []string
	}{
		{
			name:     "find interfaces 1",
			filename: "./test_classes/0x71429d7850e4421236ed7e4f58b7778fc0e0c01b8770335ba4140bcb13733e9.json",
			want: []string{
				"proxy", "proxy_l1",
			},
		}, {
			name:     "find interfaces 2",
			filename: "./test_classes/0x0702025b02d838976cf33dd7deec76b27a111d331b6093f2e7137e31c2f6ffd4.json",
			want: []string{
				"erc1155",
			},
		},
	}

	_, err := Interfaces("../../build/interfaces")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			var a abi.Abi
			err = json.NewDecoder(f).Decode(&a)
			require.NoError(t, err)

			got, err := FindInterfaces(a)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}
