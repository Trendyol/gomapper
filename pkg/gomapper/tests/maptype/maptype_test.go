package maptype

import (
	"testing"

	"github.com/Trendyol/gomapper/pkg/gomapper"
	"github.com/stretchr/testify/assert"
)

type A struct {
	MapStringX map[string]C
}

type B struct {
	MapStringX map[string]D
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
	MapStringX map[string]E
}

func Test_Map_Values(t *testing.T) {
	source := A{
		MapStringX: map[string]C{
			"key1": {Name: "struct1"},
			"key2": {Name: "struct2"},
		},
	}

	var dest B
	err := gomapper.Map(source, &dest)

	assert.Nil(t, err)
	assert.Equal(t, "struct1", dest.MapStringX["key1"].Name)
	assert.Equal(t, "struct2", dest.MapStringX["key2"].Name)
}

func Test_Empty_Map_Should_Not_Compatible(t *testing.T) {
	source := A{
		MapStringX: map[string]C{},
	}

	var dest F
	err := gomapper.Map(source, &dest)

	assert.NotNil(t, err)
}
