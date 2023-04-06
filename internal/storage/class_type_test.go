package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassType_Set(t *testing.T) {
	tests := []struct {
		name  string
		types []ClassType
		want  uint64
	}{
		{
			name:  "test 1",
			types: []ClassType{ClassTypeERC20},
			want:  1,
		}, {
			name:  "test 2",
			types: []ClassType{ClassTypeERC20, ClassTypeERC721},
			want:  3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct ClassType
			ct.Set(tt.types...)
			assert.Equal(t, ct, ClassType(tt.want))
		})
	}
}

func TestClassType_Is(t *testing.T) {
	tests := []struct {
		name  string
		ct    ClassType
		check ClassType
		want  bool
	}{
		{
			name:  "test 1",
			ct:    ClassTypeERC20,
			check: ClassTypeERC20,
			want:  true,
		}, {
			name:  "test 2",
			ct:    ClassTypeERC1155,
			check: ClassTypeERC20,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.Is(tt.check); got != tt.want {
				t.Errorf("ClassType.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
