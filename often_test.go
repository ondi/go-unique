//
//
//

package unique

import (
	"testing"

	"gotest.tools/assert"
)

var data = []struct {
	Name  string
	Count int
}{
	{"lalala", 1},
	{"bububu", 2},
	{"jujuju", 3},
}

func TestOften01(t *testing.T) {
	u := NewOften(65536)
	for _, a := range data {
		for i := 0; i < a.Count; i++ {
			u.Add(a.Name)
		}
	}
	assert.Assert(t, u.Size() == 3)

	res := u.List(len(data))
	assert.Assert(t, len(res) == len(data))
	assert.Assert(t, res[0].Value.(string) == data[2].Name)
	assert.Assert(t, res[1].Value.(string) == data[1].Name)
	assert.Assert(t, res[2].Value.(string) == data[0].Name)
}
