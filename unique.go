//
// BJKST test
//

package unique

import "sync"

type umap_t map[uint64]struct{}

type Unique_t struct {
	mx sync.Mutex
	res umap_t
	age uint
	limit int
}

// dividable by 2 ^ age
func dividable(value uint64, age uint) bool {
	return value == ((value >> age) << age)
}

func NewUnique(limit int) (self * Unique_t) {
	self = &Unique_t{}
	self.res = umap_t{}
	self.limit = limit
	return
}

func (self * Unique_t) Clear() {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.age = 0
	self.res = umap_t{}
}

func (self * Unique_t) AddUint64(value uint64) {
	self.mx.Lock()
	defer self.mx.Unlock()
	if dividable(value, self.age) == false {
		return
	}
	self.res[value] = struct{}{}
	if len(self.res) >= self.limit {
		self.age++
		for k, _ := range self.res {
			if dividable(k, self.age) == false {
				delete(self.res, k)
			}
		}
	}
}

func (self * Unique_t) SizeAge() (int, uint) {
	self.mx.Lock()
	defer self.mx.Unlock()
	return len(self.res), self.age
}

func (self * Unique_t) Count() int {
	self.mx.Lock()
	defer self.mx.Unlock()
	return len(self.res) * (1 << self.age)
}
