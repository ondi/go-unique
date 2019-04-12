//
//
//

package unique

import "testing"

var data = []string{"lalala", "bububu", "jujuju"}

func TestOften001(t * testing.T) {
	u := NewOften(65536)
	u.Add(data[0])
	u.Add(data[1])
	u.Add(data[1])
	u.Add(data[2])
	u.Add(data[2])
	u.Add(data[2])
	if u.Count() != 3 {
		t.Fatalf("Count: %v", u.Count())
	}
	var temp Results_t
	u.List(&temp, 10)
	if len(temp) != 3 {
		t.Fatalf("Result: %v", len(temp))
	}
	if temp[0].Value.(string) != data[2] {
		t.Fatalf("temp[0]")
	}
	if temp[1].Value.(string) != data[1] {
		t.Fatalf("temp[1]")
	}
	if temp[2].Value.(string) != data[0] {
		t.Fatalf("temp[2]")
	}
}
