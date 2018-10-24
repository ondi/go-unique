//
// BJKST test
//

package unique

import "sync"

type Unique_t struct {
	mx sync.Mutex
	res map[uint64]uint64
	age uint
	max_size int
}

// dividable by 2 ^ age
func dividable(value uint64, age uint) bool {
	return value == ((value >> age) << age)
}

func NewUnique(max_size int) (self * Unique_t) {
	self = &Unique_t{}
	self.res = map[uint64]uint64{}
	self.max_size = max_size
	return
}

func (self * Unique_t) Clear() {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.age = 0
	self.res = map[uint64]uint64{}
}

func (self * Unique_t) AddUint64(value uint64) {
	self.mx.Lock()
	defer self.mx.Unlock()
	if dividable(value, self.age) == false {
		return
	}
	self.res[value]++
	if len(self.res) >= self.max_size {
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

func (self * Unique_t) Count2(num uint64) (int, int) {
	self.mx.Lock()
	defer self.mx.Unlock()
	var res int
	for _, v := range self.res {
		if v <= num {
			res++
		}
	}
	return len(self.res) * (1 << self.age), res * (1 << self.age)
}
