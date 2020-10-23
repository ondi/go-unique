//
//
//

package unique

import "github.com/ondi/go-cache"

type Counter interface {
	CounterAdd(int)
	CounterGet() int
}

type Value_t struct {
	count int
}

func (self *Value_t) CounterAdd(v int) {
	self.count += v
}

func (self *Value_t) CounterGet() int {
	return self.count
}

func NewValue() Counter {
	return &Value_t{count: 1}
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

func (self *Often_t) Add(key interface{}, value func() Counter) Counter {
	it, ok := self.cc.CreateBack(key, func() interface{} { return value() })
	if !ok {
		it.Value.(Counter).CounterAdd(1)
	} else if self.cc.Size() >= self.limit {
		for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
			if it.Value.(Counter).CounterGet() == 1 {
				self.cc.Remove(it.Key)
			} else {
				it.Value.(Counter).CounterAdd(-1)
			}
		}
	}
	return it.Value.(Counter)
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
	if a.Value.(Counter).CounterGet() < b.Value.(Counter).CounterGet() {
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
