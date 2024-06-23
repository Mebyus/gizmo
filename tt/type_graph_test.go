package tt

import (
	"reflect"
	"testing"
)

func TestMergeInts(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "1 nil slices",
			a:    nil,
			b:    nil,
			want: nil,
		},
		{
			name: "2 empty slices",
			a:    []int{},
			b:    []int{},
			want: nil,
		},
		{
			name: "3 nil slice",
			a:    nil,
			b:    []int{1, 2},
			want: []int{1, 2},
		},
		{
			name: "4 nil slice",
			a:    []int{1},
			b:    nil,
			want: []int{1},
		},
		{
			name: "5 unique element",
			a:    []int{1},
			b:    []int{1},
			want: []int{1},
		},
		{
			name: "6 two single elements",
			a:    []int{1},
			b:    []int{2},
			want: []int{1, 2},
		},
		{
			name: "7 two single elements",
			a:    []int{2},
			b:    []int{1},
			want: []int{1, 2},
		},
		{
			name: "8 elements 3 + 1",
			a:    []int{2, 3, 7},
			b:    []int{1},
			want: []int{1, 2, 3, 7},
		},
		{
			name: "9 elements 1 + 3",
			a:    []int{1},
			b:    []int{2, 3, 7},
			want: []int{1, 2, 3, 7},
		},
		{
			name: "10 elements 2 + 2",
			a:    []int{1, 9},
			b:    []int{2, 8},
			want: []int{1, 2, 8, 9},
		},
		{
			name: "11 elements 3 + 5",
			a:    []int{1, 9, 10},
			b:    []int{2, 8, 12, 13, 20},
			want: []int{1, 2, 8, 9, 10, 12, 13, 20},
		},
		{
			name: "12 elements 5 + 2",
			a:    []int{5, 8, 12, 13, 20},
			b:    []int{1, 3},
			want: []int{1, 3, 5, 8, 12, 13, 20},
		},
		{
			name: "12 elements 5 + 4",
			a:    []int{5, 8, 12, 13, 20},
			b:    []int{1, 3, 5, 13},
			want: []int{1, 3, 5, 8, 12, 13, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeInts(tt.a, tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeInts() = %v, want %v", got, tt.want)
			}
		})
	}
}
