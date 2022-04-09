package app

import (
	"github.com/gin-gonic/gin"
	"testing"
	"bkapiv1/models"
	"fmt"
	"testing"
	"github.com/gin-gonic/gin"

package controllers

import (
"bkapiv1/models"
"fmt"
"testing"

"github.com/gin-gonic/gin"
)

func TestSaveProvider(t *testing.T) {
	type args struct {
		c    *gin.Context
		json ProviderProfile
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"SaveProvider",
			args{
				&gin.Context{
					//How to pass here a JSON means the below JSON ProviderProfile.
				},
				ProviderProfile{
					models.User{
						FirstName:                  "Harry",
						LastName:                   "Potter",
						FullName:                   "Harry Potter",
						CompanyName:                "TheIronNetwork",
						EmailId:                   "harry@gmail.com",
						Password:                   "vadhera123",
					},
					models.AddressStruct{},
					models.Provider{
						ProviderCategory: "IC",
						Priority:         1,
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveProvider(tt.args.c)
		})
	}
}
