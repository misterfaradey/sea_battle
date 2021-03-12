package dto

import (
	"reflect"
	"testing"
)

func Test_splitStr(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    [][2]int
		wantErr bool
	}{
		{
			name:    "1",
			in:      "1B 2C,3A 4A",
			want:    [][2]int{{0, 1}, {1, 2}, {2, 0}, {3, 0}},
			wantErr: false,
		},
		{
			name:    "2",
			in:      "1B 2C,3A 4A,",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "3",
			in:      "1B 26Z",
			want:    [][2]int{{0, 1}, {25, 25}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitStr(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitStr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
