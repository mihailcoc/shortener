package main

import (
	"testing"
)

func Test_randomInt(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomInt(tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("randomInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomString(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomString(tt.args.len); got != tt.want {
				t.Errorf("randomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
