//
// BJKST test
//

package unique

type Unique_t struct {
	res   map[uint64]struct{}
	age   uint64
	limit int
}

// dividable by 2 ^ age
func dividable(value uint64, age uint64) bool {
	return value == ((value >> age) << age)
}

func NewUnique(limit int) (self *Unique_t) {
	self = &Unique_t{}
	self.res = map[uint64]struct{}{}
	self.limit = limit
	return
}

func (self *Unique_t) Clear() {
	self.res = map[uint64]struct{}{}
	self.age = 0
}

func (self *Unique_t) AddUint64(value uint64) (added bool) {
	if added = dividable(value, self.age); added == false {
		return
	}
	self.res[value] = struct{}{}
	if len(self.res) >= self.limit {
		self.age++
		for k := range self.res {
			if dividable(k, self.age) == false {
				delete(self.res, k)
			}
		}
	}
	return
}

func (self *Unique_t) Count() int {
	return len(self.res) * (1 << self.age)
}

func (self *Unique_t) Size() int {
	return len(self.res)
}

func (self *Unique_t) Age() uint64 {
	return self.age
}
