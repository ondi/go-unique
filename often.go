//
//
//

package unique

import "sync"

type omap_t map[interface{}]uint64

type Often_t struct {
	mx sync.Mutex
	res omap_t
	limit int
}

func NewOften(limit int) (self * Often_t) {
	self = &Often_t{}
	self.res = omap_t{}
	self.limit = limit
	return
}

func (self * Often_t) Clear() {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.res = omap_t{}
}

func (self * Often_t) Add(value interface{}) {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.res[value]++
	if len(self.res) >= self.limit {
		for k, v := range self.res {
			if v == 1 {
				delete(self.res, k)
			} else {
				self.res[k] = v - 1
			}
		}
	}
}

func (self * Often_t) Count() int {
	self.mx.Lock()
	defer self.mx.Unlock()
	return len(self.res)
}
