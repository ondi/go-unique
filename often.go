//
//
//

package unique

import (
	"github.com/ondi/go-cache"
)

type Less_t = cache.Less_t[string, Counter]

// same as Range() to evict one elemts in Often_t or whole Often_t in timeline
type Evict func(f func(f func(key string, value Counter) bool))

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

func Drop(f func(f func(key string, value Counter) bool)) {}

func Less(a *cache.Value_t[string, Counter], b *cache.Value_t[string, Counter]) bool {
	return a.Value.CounterAdd(0) < b.Value.CounterAdd(0)
}

type Often_t struct {
	cc    *cache.Cache_t[string, Counter]
	evict Evict
	limit int
}

func NewOften(limit int, evict Evict) *Often_t {
	return &Often_t{
		cc:    cache.New[string, Counter](),
		evict: evict,
		limit: limit,
	}
}

func (self *Often_t) Clear() {
	self.cc = cache.New[string, Counter]()
}

func (self *Often_t) Add(key string, value func() Counter) Counter {
	it1, _ := self.cc.CreateBack(key, value)
	it1.Value.CounterAdd(1)
	if self.cc.Size() > self.limit {
		for it2 := self.cc.Front(); it2 != self.cc.End(); it2 = it2.Next() {
			if it2.Value.CounterAdd(-1) == 0 {
				self.cc.Remove(it2.Key)
				self.evict(
					func(f func(key string, value Counter) bool) {
						f(it2.Key, it2.Value)
					},
				)
			}
		}
	}
	return it1.Value
}

func (self *Often_t) Get(key string) (Counter, bool) {
	if it, ok := self.cc.Find(key); ok {
		return it.Value, true
	}
	return nil, false
}

func (self *Often_t) Size() int {
	return self.cc.Size()
}

func (self *Often_t) RangeSort(less Less_t, f func(key string, value Counter) bool) {
	self.cc.InsertionSortBack(less)
	self.Range(f)
}

func (self *Often_t) Range(f func(key string, value Counter) bool) {
	for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
		if f(it.Key, it.Value) == false {
			return
		}
	}
}
