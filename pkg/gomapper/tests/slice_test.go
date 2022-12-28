package tests

import (
	"testing"

	"github.com/Trendyol/gomapper/pkg/gomapper"
	"github.com/stretchr/testify/assert"
)

type RootX struct {
	Flavor *FlavorX
}

type FlavorX struct {
	Type  string
	Roles *[]RoleX
}

type RoleX struct {
	Size  *int
	Count *int
}

type RootY struct {
	Flavor *FlavorY
}

type FlavorY struct {
	Type  string
	Roles *[]RoleY
}

type RoleY struct {
	Size  *int
	Count *int
}

func Test_Slice(t *testing.T) {
	size := 50
	count := 5

	roleSlice := []RoleX{
		{
			Size:  &size,
			Count: &count,
		},
	}

	source := FlavorX{
		Type:  "small",
		Roles: &roleSlice,
	}

	var dest FlavorY
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.NotNil(t, dest.Roles)
}

func Test_Slice_Nil(t *testing.T) {
	source := FlavorX{
		Type:  "small",
		Roles: nil,
	}

	var dest FlavorY
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.Nil(t, dest.Roles)
}
