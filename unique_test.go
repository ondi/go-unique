//
//
//

package unique

import "testing"

func TestUnique001(t * testing.T) {
	u := NewUnique(65536)
	u.AddUint64(1)
	if u.Count() != 1 {
		t.Fatalf("1: %v", u.Count())
	}
}
