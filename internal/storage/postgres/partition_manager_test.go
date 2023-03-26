package postgres

import (
	"reflect"
	"testing"
	"time"
)

func Test_quarterOf(t *testing.T) {
	tests := []struct {
		name  string
		month time.Month
		want  int
	}{
		{
			name:  "Jan",
			month: time.January,
			want:  1,
		}, {
			name:  "Feb",
			month: time.February,
			want:  1,
		}, {
			name:  "Mar",
			month: time.March,
			want:  1,
		}, {
			name:  "Apr",
			month: time.April,
			want:  2,
		}, {
			name:  "May",
			month: time.May,
			want:  2,
		}, {
			name:  "Jun",
			month: time.June,
			want:  2,
		}, {
			name:  "Jul",
			month: time.July,
			want:  3,
		}, {
			name:  "Aug",
			month: time.August,
			want:  3,
		}, {
			name:  "Sep",
			month: time.September,
			want:  3,
		}, {
			name:  "Oct",
			month: time.October,
			want:  4,
		}, {
			name:  "Nov",
			month: time.November,
			want:  4,
		}, {
			name:  "Dec",
			month: time.December,
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quarterOf(tt.month); got != tt.want {
				t.Errorf("quarterOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_quarterBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		current time.Time
		want    time.Time
		want1   time.Time
		wantErr bool
	}{
		{
			name:    "test 1",
			current: time.Date(2022, time.January, 12, 2, 2, 2, 2, time.UTC),
			want:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			want1:   time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
		}, {
			name:    "test 2",
			current: time.Date(2022, time.May, 12, 2, 2, 2, 2, time.UTC),
			want:    time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
			want1:   time.Date(2022, time.July, 1, 0, 0, 0, 0, time.UTC),
		}, {
			name:    "test 3",
			current: time.Date(2022, time.August, 12, 2, 2, 2, 2, time.UTC),
			want:    time.Date(2022, time.July, 1, 0, 0, 0, 0, time.UTC),
			want1:   time.Date(2022, time.October, 1, 0, 0, 0, 0, time.UTC),
		}, {
			name:    "test 4",
			current: time.Date(2022, time.October, 12, 2, 2, 2, 2, time.UTC),
			want:    time.Date(2022, time.October, 1, 0, 0, 0, 0, time.UTC),
			want1:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := quarterBoundaries(tt.current)
			if (err != nil) != tt.wantErr {
				t.Errorf("quarterBoundaries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("quarterBoundaries() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("quarterBoundaries() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
