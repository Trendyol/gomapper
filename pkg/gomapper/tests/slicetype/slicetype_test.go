package slicetype

import (
	"testing"

	"github.com/Trendyol/gomapper/pkg/gomapper"
	"github.com/stretchr/testify/assert"
)

type A struct {
	X *C
}

type B struct {
	X *D
}

type C struct {
	TestString string
	Slice      *[]E
}

type CNonPtrESlice struct {
	TestString string
	Slice      []E
}

type D struct {
	TestString string
	Slice      *[]F
}

type E struct {
	SizeInt  *int
	CountInt *int
}

type F struct {
	SizeInt  *int
	CountInt *int
}

func Test_From_Root(t *testing.T) {
	size := 50
	count := 5

	source := A{X: &C{
		TestString: "test",
		Slice: &[]E{
			{
				SizeInt:  &size,
				CountInt: &count,
			},
		},
	}}

	var dest B
	err := gomapper.Map(source, &dest)

	assert.Nil(t, err)
	assert.NotNil(t, dest.X.Slice)
}

func Test_Slice_When_Dest_Ptr(t *testing.T) {
	size := 50
	count := 5

	source := C{
		TestString: "test",
		Slice: &[]E{
			{
				SizeInt:  &size,
				CountInt: &count,
			},
		},
	}

	var dest D
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.NotNil(t, dest.Slice)
}

func Test_Slice_When_Dest_NonPtr(t *testing.T) {
	size := 50
	count := 5

	source := C{
		TestString: "test",
		Slice: &[]E{
			{
				SizeInt:  &size,
				CountInt: &count,
			},
		}}

	var dest CNonPtrESlice
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.NotNil(t, dest.Slice)
}

func Test_Slice_Nil_When_Dest_Ptr(t *testing.T) {
	source := C{
		TestString: "test",
		Slice:      nil,
	}

	var dest D
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.Nil(t, dest.Slice)
}

func Test_Slice_Nil_When_Dest_NonPtr(t *testing.T) {
	source := C{
		TestString: "test",
		Slice:      nil,
	}

	var dest CNonPtrESlice
	err := gomapper.Map(&source, &dest)

	assert.Nil(t, err)
	assert.Nil(t, dest.Slice)
}
