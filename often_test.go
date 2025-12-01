//
//
//

package unique

import (
	"testing"

	"gotest.tools/assert"
)

type Value_t struct {
	count int64
}

func (self *Value_t) CounterAdd(a int64) {
	self.count += a
}

func (self *Value_t) CounterGet() int64 {
	return self.count
}

var data = []struct {
	Name  string
	Count int
}{
	{"lalala", 1},
	{"bububu", 2},
	{"jujuju", 3},
}

type Result_t struct {
	Key   string
	Value *Value_t
}

type ResultList_t []Result_t

func (self *ResultList_t) Add(key string, value *Value_t) bool {
	*self = append(*self, Result_t{key, value})
	return true
}

func TestOften01(t *testing.T) {
	u := NewOften(65536, Drop[string, *Value_t])
	for _, a := range data {
		for i := 0; i < a.Count; i++ {
			u.Add(a.Name, func(p **Value_t) { *p = &Value_t{} }, func(p **Value_t) {})
		}
	}
	assert.Assert(t, u.Size() == 3)

	var res ResultList_t
	u.RangeSort(Less[string, *Value_t], res.Add)
	assert.Assert(t, len(res) == len(data))
	assert.Assert(t, res[0].Key == data[2].Name)
	assert.Assert(t, res[1].Key == data[1].Name)
	assert.Assert(t, res[2].Key == data[0].Name)
}
