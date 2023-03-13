package mapkind

import (
	"testing"

	"github.com/Trendyol/gomapper/pkg/gomapper"
	"github.com/stretchr/testify/assert"
)

type A struct {
	MapX map[string]C
}

type B struct {
	MapX map[string]D
}

type C struct {
	Name string
}

type D struct {
	Name string
}

type E struct {
	Name int
}

type F struct {
	MapX map[string]E
}

type G struct {
	MapX map[int]E
}

func Test_Map_Values_Equality(t *testing.T) {
	source := A{
		MapX: map[string]C{
			"key1": {Name: "struct1"},
			"key2": {Name: "struct2"},
		},
	}

	var dest B
	err := gomapper.Map(source, &dest)

	assert.Nil(t, err)
	assert.Equal(t, "struct1", dest.MapX["key1"].Name)
	assert.Equal(t, "struct2", dest.MapX["key2"].Name)
}

func Test_Map_Element_Types_Should_Not_Compatible(t *testing.T) {
	source := A{
		MapX: map[string]C{},
	}

	var dest F
	err := gomapper.Map(source, &dest)

	assert.NotNil(t, err)
}

func Test_Map_Key_Types_Should_Not_Compatible(t *testing.T) {
	source := F{
		MapX: map[string]E{},
	}

	var dest G
	err := gomapper.Map(source, &dest)

	assert.NotNil(t, err)
}
