//
//
//

package unique

import "testing"

func TestOften001(t * testing.T) {
	u := NewOften(65536)
	u.Add("lalala")
	if u.Count() != 1 {
		t.Fatalf("1: %v", u.Count())
	}
}
