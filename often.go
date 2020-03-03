//
//
//

package unique

import "github.com/ondi/go-cache"

type Value interface {
	CountAdd(int)
	CountGet() int
}

type Value_t struct {
	count int
}

func (self *Value_t) CountAdd(v int) {
	self.count += v
}

func (self *Value_t) CountGet() int {
	return self.count
}

func NewValue() Value {
	return &Value_t{count: 1}
}

type Result_t struct {
	Key   interface{}
	Value Value
}

type Often_t struct {
	cc    *cache.Cache_t
	limit int
}

func NewOften(limit int) (self *Often_t) {
	self = &Often_t{}
	self.cc = cache.New()
	self.limit = limit
	return
}

func (self *Often_t) Clear() {
	self.cc = cache.New()
}

func (self *Often_t) Add(key interface{}, value func() Value) Value {
	it, ok := self.cc.CreateBack(key, func() interface{} { return value() })
	if !ok {
		it.Value().(Value).CountAdd(1)
	} else if self.cc.Size() >= self.limit {
		for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
			if it.Value().(Value).CountGet() == 1 {
				self.cc.Remove(it.Key())
			} else {
				it.Value().(Value).CountAdd(-1)
			}
		}
	}
	return it.Value().(Value)
}

func (self *Often_t) Size() int {
	return self.cc.Size()
}

type Less_t struct{}

func (Less_t) Less(a *cache.Value_t, b *cache.Value_t) bool {
	if a.Value().(Value).CountGet() < b.Value().(Value).CountGet() {
		return true
	}
	return false
}

func (self *Often_t) List(limit int) (res []Result_t) {
	self.cc.InsertionSortBack(Less_t{})
	for it := self.cc.Front(); it != self.cc.End() && limit > 0; it = it.Next() {
		res = append(res, Result_t{Key: it.Key(), Value: it.Value().(Value)})
		limit--
	}
	return
}
