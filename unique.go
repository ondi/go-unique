//
// BJKST test
//

package unique

type Unique_t struct {
	res   map[uint64]struct{}
	age   uint
	limit int
}

// dividable by 2 ^ age
func dividable(value uint64, age uint) bool {
	return value == ((value >> age) << age)
}

func NewUnique(limit int) (self *Unique_t) {
	self = &Unique_t{}
	self.res = map[uint64]struct{}{}
	self.limit = limit
	return
}

func (self *Unique_t) Clear() {
	self.age = 0
	self.res = map[uint64]struct{}{}
}

func (self *Unique_t) AddUint64(value uint64) {
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

func (self *Unique_t) SizeAge() (int, uint) {
	return len(self.res), self.age
}

func (self *Unique_t) Count() int {
	return len(self.res) * (1 << self.age)
}
