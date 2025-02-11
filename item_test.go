package lik

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItem(t *testing.T) {
	data := `{"str":"str", "yes":true, "no":false, "digit":{"dva":2,"nepolt":-1.5}, "array":[11, "12", 13]}`
	item := ItemFromString(data)
	if item == nil {
		t.Errorf("ERROR of ItemFromString(): noitem")
	} else if set := item.ToSet(); set == nil {
		t.Errorf("ERROR of ItemFromString(): noset")
	} else {
		assert.Equal(t, "str", set.GetString("str"))
		assert.Equal(t, false, set.GetBool("no"))
		assert.Equal(t, true, set.GetBool("yes"))
		assert.Equal(t, 2, int(set.GetInt("digit/dva")))
		assert.Equal(t, -1.5, set.GetFloat("digit/nepolt"))
		assert.Equal(t, 11, int(set.GetInt("array/0")))
		assert.Equal(t, 12, int(set.GetInt("array/1")))
		assert.Equal(t, "12", set.GetString("array/1"))
		assert.Equal(t, 13, int(set.GetInt("array/2")))
	}
}
