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

// same as RangeRaw() to evict one elemts in Often_t or whole Often_t in timeline
type Evict func(f func(f func(key string, value Counter) bool))

type Value_t struct {
	count int64
}

func (self *Value_t) CounterAdd(a int64) {
	self.count += a
}

func (self *Value_t) CounterGet() int64 {
	return self.count
}

func Drop(f func(f func(key string, value Counter) bool)) {}

func Less(a *cache.Value_t[string, Counter], b *cache.Value_t[string, Counter]) bool {
	if any(a.Value).(Counter).CounterGet() < any(b.Value).(Counter).CounterGet() {
		return true
	}
	return false
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

func (self *Often_t) Add(key string, value func() Counter) (res Counter) {
	it, _ := self.cc.CreateBack(key, value)
	res = it.Value.(Counter)
	res.CounterAdd(1)
	if self.cc.Size() > self.limit {
		for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
			if it.Value.(Counter).CounterGet() == 1 {
				self.cc.Remove(it.Key)
				self.evict(
					func(f func(key string, value Counter) bool) {
						f(it.Key, it.Value.(Counter))
					},
				)
			} else {
				it.Value.(Counter).CounterAdd(-1)
			}
		}
	}
	return
}

func (self *Often_t) Get(key string) (Counter, bool) {
	if it, ok := self.cc.Find(key); ok {
		return it.Value.(Counter), true
	}
	return nil, false
}

func (self *Often_t) Size() int {
	return self.cc.Size()
}

func (self *Often_t) Range(less cache.Less[string, Counter], f func(key string, value Counter) bool) {
	self.cc.InsertionSortBack(less)
	self.RangeRaw(f)
}

func (self *Often_t) RangeRaw(f func(key string, value Counter) bool) {
	for it := self.cc.Front(); it != self.cc.End(); it = it.Next() {
		if f(it.Key, it.Value.(Counter)) == false {
			return
		}
	}
}
