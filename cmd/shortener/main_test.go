package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_handlerGet(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerGet(tt.args.g)
		})
	}
}

func Test_handlerPost(t *testing.T) {
	type args struct {
		g *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerPost(tt.args.g)
		})
	}
}
