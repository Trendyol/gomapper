package tests

import (
	"testing"

	"github.com/Trendyol/gomapper/pkg/gomapper"
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

type E struct {
	Name Z
}

type F struct {
	Name *string
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

	err := gomapper.Map(source, dest)
	assert.NotNil(t, err)

	err = gomapper.Map(nil, dest)
	assert.NotNil(t, err)
}

func Test_Dest_Must_Not_Nil(t *testing.T) {
	dest := (*X)(nil)

	err := gomapper.Map(X{}, dest)
	assert.NotNil(t, err)

	err = gomapper.Map(X{}, nil)
	assert.NotNil(t, err)
}

func Test_Dest_Must_Not_Nil_Panic(t *testing.T) {
	dest := (*X)(nil)

	assert.Panics(t, func() { gomapper.MapP(X{}, dest) })
	assert.Panics(t, func() { gomapper.MapP(X{}, nil) })
}

func Test_Dest_Must_Be_Pointer(t *testing.T) {
	err := gomapper.Map(X{}, X{})
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

	if err := gomapper.Map(source, dest); err != nil {
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

	if err := gomapper.Map(source, dest); err != nil {
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

	err := gomapper.Map(source, dest, &gomapper.Option{Exact: true})
	assert.NotNil(t, err)
}

func Test_Y_To_Z_Map_Loose(t *testing.T) {
	source := Y{
		In: B{C: C{address: "istanbul"}},
	}

	dest := &Z{}

	if err := gomapper.Map(source, dest); err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, source.In.C.address, dest.In.C.address)
}

func Test_Mapping_Source_String_Field_To_Dest_Struct_Field_Should_Return_Error(t *testing.T) {
	source := A{
		Name: "istanbul",
		zone: []string{},
	}

	dest := &E{}

	err := gomapper.Map(source, dest)

	assert.NotNil(t, err)
}

func Test_Mapping_Nil_Pointer_Source_String_Field_To_Dest_Struct_Field_Should_Return_Error(t *testing.T) {
	source := F{
		Name: nil,
	}

	dest := &E{}

	err := gomapper.Map(source, dest)

	assert.NotNil(t, err)
}
