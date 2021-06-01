//
//
//

package unique

import "github.com/ondi/go-cache"

type Counter interface {
	CounterAdd(int64) int64
}

type Value_t struct {
	count int64
}

func (self *Value_t) CounterAdd(a int64) int64 {
	self.count += a
	return self.count
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

func (self *Often_t) Add(key interface{}, value func() Counter) (res Counter) {
	it, _ := self.cc.CreateBack(key, func() interface{} { return value() })
	res = it.Value.(Counter)
	res.CounterAdd(1)
	if self.cc.Size() >= self.limit {
		for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
			if it.Value.(Counter).CounterAdd(0) == 1 {
				self.cc.Remove(it.Key)
			} else {
				it.Value.(Counter).CounterAdd(-1)
			}
		}
	}
	return
}

func (self *Often_t) Get(key interface{}) (Counter, bool) {
	if it, ok := self.cc.Find(key); ok {
		return it.Value.(Counter), true
	}
	return nil, false
}

func (self *Often_t) Size() int {
	return self.cc.Size()
}

type Less_t struct{}

func (Less_t) Less(a *cache.Value_t, b *cache.Value_t) bool {
	if a.Value.(Counter).CounterAdd(0) < b.Value.(Counter).CounterAdd(0) {
		return true
	}
	return false
}

func (self *Often_t) Range(less cache.IsLess, f func(key interface{}, value Counter) bool) {
	self.cc.InsertionSortBack(less)
	for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
		if f(it.Key, it.Value.(Counter)) == false {
			return
		}
	}
}
