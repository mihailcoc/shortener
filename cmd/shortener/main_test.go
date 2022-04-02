package main

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_viewHandler(t *testing.T) {
	tests := []struct {
		name string
		want *gin.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := viewHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("viewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
