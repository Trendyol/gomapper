package gomapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	Name string
	zone []string
}

type X struct {
	Ui64 uint64
	i64  int64
	A    A
}

type XPointerField struct {
	Ui64 uint64
	i64  int64
	A    *A
}

type B struct {
	C C
}

type D struct {
	C *C
}

type C struct {
	address string
}

type Y struct {
	In B
}

type Z struct {
	In *D
}

func Test_Source_Must_Not_Nil(t *testing.T) {
	source := (*X)(nil)

	dest := &X{}

	err := MapLoose(source, dest)
	assert.NotNil(t, err)

	err = MapLoose(nil, dest)
	assert.NotNil(t, err)
}

func Test_Dest_Must_Not_Nil(t *testing.T) {
	dest := (*X)(nil)

	err := MapLoose(X{}, dest)
	assert.NotNil(t, err)

	err = MapLoose(X{}, nil)
	assert.NotNil(t, err)
}

func Test_Dest_Must_Be_Pointer(t *testing.T) {
	err := MapLoose(X{}, X{})
	assert.NotNil(t, err)
}

func Test_X_To_X_Map_Loose(t *testing.T) {
	source := X{
		Ui64: 123,
		i64:  321,
		A: A{
			Name: "Abc",
			zone: []string{"a", "b", "c", "d"},
		},
	}

	dest := &X{}

	if err := MapLoose(source, dest); err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, source.Ui64, dest.Ui64)
	assert.Equal(t, source.i64, dest.i64)
	assert.Equal(t, source.A.Name, dest.A.Name)
	assert.Equal(t, source.A.zone, dest.A.zone)
}

func Test_X_To_XPointerField_Map_Loose(t *testing.T) {
	source := X{
		Ui64: 123,
		i64:  321,
		A: A{
			Name: "Abc",
			zone: []string{"a", "b", "c", "d"},
		},
	}

	dest := &XPointerField{}

	if err := MapLoose(source, dest); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, source.Ui64, dest.Ui64)

	// Can't map private fields when types are different.
	assert.Equal(t, int64(0), dest.i64)

	assert.Equal(t, source.A.Name, dest.A.Name)
	assert.Equal(t, source.A.zone, dest.A.zone)
}

func Test_X_To_XPointerField_Map(t *testing.T) {
	source := X{
		Ui64: 123,
		i64:  321,
		A: A{
			Name: "Abc",
			zone: []string{"a", "b", "c", "d"},
		},
	}

	dest := &XPointerField{}

	err := Map(source, dest)
	assert.NotNil(t, err)
}

func Test_Y_To_Z_Map_Loose(t *testing.T) {
	source := Y{
		In: B{C: C{address: "istanbul"}},
	}

	dest := &Z{}

	if err := MapLoose(source, dest); err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, source.In.C.address, dest.In.C.address)
}
