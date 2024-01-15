package utils

import (
	"reflect"
	"strconv"
	"testing"
)

type typeStringer struct {
	i int
}

type typeNotStringer struct {
	i int
}

func (s typeStringer) String() string {
	return strconv.Itoa(s.i)
}

func TestKeyToString(t *testing.T) {
	t.Parallel()

	type args struct {
		key interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{key: "string"},
			want: "string",
		},
		{
			name: "stringer",
			args: args{key: typeStringer{i: 81}},
			want: "81",
		},
		{
			name: "bytes",
			args: args{key: []byte("byte")},
			want: "byte",
		},
		{
			name: "int",
			args: args{key: 8},
			want: "8",
		},
		{
			name: "int8",
			args: args{key: int8(8)},
			want: "8",
		},
		{
			name: "int16",
			args: args{key: int16(8)},
			want: "8",
		},
		{
			name: "int32",
			args: args{key: int32(8)},
			want: "8",
		},
		{
			name: "int64",
			args: args{key: int64(8)},
			want: "8",
		},
		{
			name: "uint",
			args: args{key: uint(8)},
			want: "8",
		},
		{
			name: "uint8",
			args: args{key: uint8(8)},
			want: "8",
		},
		{
			name: "uint16",
			args: args{key: uint16(8)},
			want: "8",
		},
		{
			name: "uint32",
			args: args{key: uint32(8)},
			want: "8",
		},
		{
			name: "uint64",
			args: args{key: uint64(8)},
			want: "8",
		},
		{
			name: "float32",
			args: args{key: float32(8.8)},
			want: "8.8",
		},
		{
			name: "float64",
			args: args{key: 8.8},
			want: "8.8",
		},
		{
			name: "notStringer",
			args: args{key: typeNotStringer{i: 81}},
			want: "{81}",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := KeyToString(tt.args.key); got != tt.want {
				t.Errorf("KeyToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeysToString(t *testing.T) {
	t.Parallel()

	want := []string{"123", "123"}

	t.Run("strings", func(t *testing.T) {
		t.Parallel()
		checkSlice(t, []string{"123", "123"}, want)
	})

	t.Run("int8", func(t *testing.T) {
		t.Parallel()
		checkSlice(t, []int8{123, 123}, want)
	})

	t.Run("uint8", func(t *testing.T) {
		t.Parallel()
		checkSlice(t, []uint8{123, 123}, want)
	})
}
func checkSlice[T any](t *testing.T, in []T, want []string) {
	if got := KeysToString(in); !reflect.DeepEqual(got, want) {
		t.Errorf("KeysToString() = %v, want %v", got, want)
	}
}
