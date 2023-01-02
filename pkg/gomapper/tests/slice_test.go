package tests

import (
	"testing"

	"github.com/Trendyol/gomapper/pkg/gomapper"
	"github.com/stretchr/testify/assert"
)

type FlavorX struct {
	Type  string
	Roles *[]RoleX
}

type FlavorY struct {
	Type  string
	Roles *[]RoleY
}

type FlavorWithNonPointerRoleSlice struct {
	Type  string
	Roles []RoleX
}

type RootX struct {
	Flavor *FlavorX
}

type RootY struct {
	Flavor *FlavorY
}

type RoleX struct {
	Size  *int
	Count *int
}

type RoleY struct {
	Size  *int
	Count *int
}

func Test_Slice_When_Dest_Ptr(t *testing.T) {
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

func Test_Slice_When_Dest_NonPtr(t *testing.T) {
	size := 50
	count := 5

	source := FlavorX{
		Type: "small",
		Roles: &[]RoleX{
			{
				Size:  &size,
				Count: &count,
			},
		}}

	var dest FlavorWithNonPointerRoleSlice
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.NotNil(t, dest.Roles)
}

func Test_Slice_Nil_When_Dest_Ptr(t *testing.T) {
	source := FlavorX{
		Type:  "small",
		Roles: nil,
	}

	var dest FlavorY
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.Nil(t, dest.Roles)
}

func Test_Slice_Nil_When_Dest_NonPtr(t *testing.T) {
	source := FlavorX{
		Type:  "small",
		Roles: nil,
	}

	var dest FlavorWithNonPointerRoleSlice
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.Nil(t, dest.Roles)
}
