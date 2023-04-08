package filters

import (
	"strings"
	"testing"
)

func Test_findPathInMap(t *testing.T) {
	type args struct {
		path string
		m    map[string]any
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "test 1",
			args: args{
				path: "key",
				m:    map[string]any{"key": "test"},
			},
			want:  "test",
			want1: true,
		}, {
			name: "test 2",
			args: args{
				path: "key",
				m:    map[string]any{"key1": "test"},
			},
			want:  "",
			want1: false,
		}, {
			name: "test 3",
			args: args{
				path: "key.subkey.subkey2",
				m: map[string]any{
					"key": map[string]any{
						"subkey": map[string]any{
							"subkey2": "test",
						},
					},
				},
			},
			want:  "test",
			want1: true,
		}, {
			name: "test 4",
			args: args{
				path: "key.subkey.subkey3",
				m: map[string]any{
					"key": map[string]any{
						"subkey": map[string]any{
							"subkey2": "test",
						},
					},
				},
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := strings.Split(tt.args.path, ".")
			got, got1 := findPathInMap(path, tt.args.m)
			if got != tt.want {
				t.Errorf("findPathInMap() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findPathInMap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
