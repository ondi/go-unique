//
//
//

package unique

import "github.com/ondi/go-cache"

type Result_t struct {
	Value interface{}
	Count int
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

func (self *Often_t) Add(value interface{}) {
	if it, ok := self.cc.CreateBack(value, func() interface{} { return 1 }); !ok {
		it.Update(it.Value().(int) + 1)
	} else if self.cc.Size() >= self.limit {
		for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
			if it.Value().(int) == 1 {
				self.cc.Remove(it.Key())
			} else {
				it.Update(it.Value().(int) - 1)
			}
		}
	}
}

func (self *Often_t) Size() int {
	return self.cc.Size()
}

type Less_t struct{}

func (Less_t) Less(a *cache.Value_t, b *cache.Value_t) bool {
	if a.Value().(int) < b.Value().(int) {
		return true
	}
	return false
}

func (self *Often_t) List(limit int) (res []Result_t) {
	self.cc.InsertionSortBack(Less_t{})
	for it := self.cc.Front(); it != self.cc.End() && limit > 0; it = it.Next() {
		res = append(res, Result_t{it.Key(), it.Value().(int)})
		limit--
	}
	return
}
