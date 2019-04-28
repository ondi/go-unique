//
//
//

package unique

import "sync"

import "github.com/ondi/go-cache"

type Result_t struct {
	Value interface{}
	Count int
}

type Append interface {
	Append(Result_t)
}

type Results_t []Result_t

func (self * Results_t) Append(res Result_t) {
	*self = append(*self, res)
}

type Often_t struct {
	mx sync.Mutex
	res * cache.Cache_t
	limit int
}

func NewOften(limit int) (self * Often_t) {
	self = &Often_t{}
	self.res = cache.New()
	self.limit = limit
	return
}

func (self * Often_t) Clear() {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.res = cache.New()
}

func (self * Often_t) Add(value interface{}) {
	self.mx.Lock()
	defer self.mx.Unlock()
	if it, ok := self.res.CreateBack(value, 1); !ok {
		it.Update(it.Value().(int) + 1)
	} else if self.res.Size() >= self.limit {
		for it := self.res.Front(); it != self.res.End(); it = it.Next() {
			if it.Value().(int) == 1 {
				self.res.Remove(it.Key())
			} else {
				it.Update(it.Value().(int) - 1)
			}
		}
	}
}

func (self * Often_t) Count() int {
	self.mx.Lock()
	defer self.mx.Unlock()
	return self.res.Size()
}

type Less_t struct {}

func (Less_t) Less(a * cache.Value_t, b * cache.Value_t) bool {
	if a.Value().(int) < b.Value().(int) {
		return true
	}
	return false
}

func (self * Often_t) List(a Append, limit int) {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.res.InsertionSortBack(Less_t{})
	for it := self.res.Front(); it != self.res.End() && limit > 0; it = it.Next() {
		a.Append(Result_t{it.Key(), it.Value().(int)})
		limit--
	}
}
