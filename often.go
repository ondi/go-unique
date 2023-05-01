//
//
//

package unique

import (
	"github.com/ondi/go-cache"
)

type Counter interface {
	CounterAdd(int64)
	CounterGet() int64
}

type Value_t struct {
	count int64
}

func (self *Value_t) CounterAdd(a int64) {
	self.count += a
}

func (self *Value_t) CounterGet() int64 {
	return self.count
}

type Evict_t[Mapped_t Counter] func(key string, value Mapped_t)

func Drop[Mapped_t Counter](key string, value Mapped_t) {}

func Less[Mapped_t Counter](a *cache.Value_t[string, Mapped_t], b *cache.Value_t[string, Mapped_t]) bool {
	return a.Value.CounterGet() < b.Value.CounterGet()
}

type Often_t[Mapped_t Counter] struct {
	cc    *cache.Cache_t[string, Mapped_t]
	evict Evict_t[Mapped_t]
	limit int
}

func NewOften[Mapped_t Counter](limit int, evict Evict_t[Mapped_t]) *Often_t[Mapped_t] {
	return &Often_t[Mapped_t]{
		cc:    cache.New[string, Mapped_t](),
		evict: evict,
		limit: limit,
	}
}

func (self *Often_t[Mapped_t]) Clear() {
	self.cc = cache.New[string, Mapped_t]()
}

func (self *Often_t[Mapped_t]) Add(key string, value func(*Mapped_t)) (res Mapped_t, ok bool) {
	it1, ok := self.cc.CreateBack(key, value, func(*Mapped_t) {})
	it1.Value.CounterAdd(1)
	if self.cc.Size() > self.limit {
		for it2 := self.cc.Front(); it2 != self.cc.End(); it2 = it2.Next() {
			if it2.Value.CounterAdd(-1); it2.Value.CounterGet() <= 0 {
				self.cc.Remove(it2.Key)
				self.evict(it2.Key, it2.Value)
				break
			}
		}
	}
	res = it1.Value
	return
}

func (self *Often_t[Mapped_t]) Get(key string) (out Mapped_t, ok bool) {
	it, ok := self.cc.Find(key)
	if ok {
		out = it.Value
	}
	return
}

func (self *Often_t[Mapped_t]) Size() int {
	return self.cc.Size()
}

func (self *Often_t[Mapped_t]) RangeSort(less cache.Less_t[string, Mapped_t], f func(key string, value Mapped_t) bool) {
	self.cc.InsertionSortBack(less)
	self.Range(f)
}

func (self *Often_t[Mapped_t]) Range(f func(key string, value Mapped_t) bool) {
	for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
		if f(it.Key, it.Value) == false {
			return
		}
	}
}
