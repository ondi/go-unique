//
// BJKST test
//

package unique

import "sync"

type Unique_t struct {
	mx sync.Mutex
	res map[uint64]struct{}
	age uint
	samples int
}

func NewUnique(samples int) (self * Unique_t) {
	self = &Unique_t{}
	self.res = map[uint64]struct{}{}
	self.samples = samples
	return
}

func (self * Unique_t) Clear() {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.age = 0
	self.res = map[uint64]struct{}{}
}

// dividable by 2 ^ age
func (self * Unique_t) dividable(value uint64) bool {
	return value == ((value >> self.age) << self.age)
}

func (self * Unique_t) add(value uint64) {
	if self.dividable(value) == false {
		return
	}
	self.mx.Lock()
	defer self.mx.Unlock()
	self.res[value] = struct{}{}
	if len(self.res) >= self.samples {
		self.age++
		for k, _ := range self.res {
			if self.dividable(k) == false {
				delete(self.res, k)
			}
		}
	}
}

func (self * Unique_t) AddUint64(value uint64) {
	// value should be hashed
	self.add(value)
}

func (self * Unique_t) Size() int {
	self.mx.Lock()
	defer self.mx.Unlock()
	return len(self.res) * (1 << self.age)
}

func (self * Unique_t) Size2() (int, uint) {
	self.mx.Lock()
	defer self.mx.Unlock()
	return len(self.res), self.age
}
