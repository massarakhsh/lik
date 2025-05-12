package lik

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItem(t *testing.T) {
	data := `{"str":"str", "yes":true, "no":false, "net":null, "digit":{"dva":2,"nepolt":-1.5}, "array":[11, "12", null, 13]}`
	item := ItemFromString(data)
	if item == nil {
		t.Errorf("ERROR of ItemFromString(): noitem")
	} else if set := item.ToSet(); set == nil {
		t.Errorf("ERROR of ItemFromString(): noset")
	} else {
		set.SetValue("append/one", "one")
		set.SetValue("append/two", "two")
		set.SetValue("=/long/long", "longlong")
		//fmt.Println(set.Format(""))
		assert.Equal(t, "str", set.GetString("str"))
		assert.Equal(t, false, set.GetBool("no"))
		assert.Equal(t, true, set.GetBool("yes"))
		assert.Nil(t, set.GetItem("net"))
		assert.Equal(t, 2, int(set.GetInt("digit/dva")))
		assert.Equal(t, -1.5, set.GetFloat("digit/nepolt"))
		assert.Equal(t, 3, int(set.GetList("array").Count()))
		assert.Equal(t, 11, int(set.GetInt("array/0")))
		assert.Equal(t, 12, int(set.GetInt("array/1")))
		assert.Equal(t, "12", set.GetString("array/1"))
		assert.Equal(t, 13, int(set.GetInt("array/2")))
		assert.Equal(t, "one", set.GetString("append/one"))
		assert.Equal(t, "two", set.GetString("append/two"))
		assert.Nil(t, set.GetItem("="))
		assert.Nil(t, set.GetItem("long"))
		assert.Nil(t, set.GetItem("/long"))
		assert.Nil(t, set.GetItem("=/long"))
		assert.NotNil(t, set.GetItem("=/long/long"))
		assert.Equal(t, "longlong", set.GetString("=/long/long"))
	}
}
