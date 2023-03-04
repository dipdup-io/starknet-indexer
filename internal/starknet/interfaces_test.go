package starknet

import (
	"testing"

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
			wantLen: 5,
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
