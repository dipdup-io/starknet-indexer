package starknet

import (
	"os"
	"reflect"
	"testing"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
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
			wantLen: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Interfaces(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interfaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, tt.wantLen, "interfaces count")
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
			name:     "find interfaces",
			filename: "./test_classes/0x71429d7850e4421236ed7e4f58b7778fc0e0c01b8770335ba4140bcb13733e9.json",
			want: []string{
				"proxy",
			},
		}, {
			name:     "find interfaces",
			filename: "./test_classes/0x0702025b02d838976cf33dd7deec76b27a111d331b6093f2e7137e31c2f6ffd4.json",
			want: []string{
				"erc1155",
			},
		},
	}

	if _, err := Interfaces("../../build/interfaces"); err != nil {
		t.Errorf("can't load interafces: %s", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.filename)
			if err != nil {
				t.Errorf("can't open file: %s", err)
				return
			}
			defer f.Close()

			var a abi.Abi
			if err := json.NewDecoder(f).Decode(&a); err != nil {
				t.Errorf("can't decode abi^ %s", err)
				return
			}

			got, err := FindInterfaces(a)
			if err != nil {
				t.Errorf("can't find interface: %s", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("invalid interfaces set")
			}
		})
	}
}
