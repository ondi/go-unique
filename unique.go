//
// BJKST test
//

package unique

type Unique_t struct {
	samples map[uint64]struct{}
	age     uint64
	limit   int
}

// dividable by 2 ^ age
func dividable(value uint64, age uint64) bool {
	return value == ((value >> age) << age)
}

func NewUnique(limit int) (self *Unique_t) {
	self = &Unique_t{}
	self.samples = map[uint64]struct{}{}
	self.limit = limit
	return
}

func (self *Unique_t) Clear() {
	self.samples = map[uint64]struct{}{}
	self.age = 0
}

func (self *Unique_t) AddUint64(value uint64) (added bool) {
	if added = dividable(value, self.age); !added {
		return
	}
	if _, added = self.samples[value]; added {
		return false
	}
	self.samples[value] = struct{}{}
	if len(self.samples) >= self.limit {
		self.age++
		for k := range self.samples {
			if !dividable(k, self.age) {
				delete(self.samples, k)
			}
		}
	}
	return true
}

func (self *Unique_t) Count() int {
	return len(self.samples) * (1 << self.age)
}

func (self *Unique_t) Size() int {
	return len(self.samples)
}

func (self *Unique_t) Age() uint64 {
	return self.age
}
