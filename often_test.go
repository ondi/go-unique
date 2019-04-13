//
//
//

package unique

import "testing"

var data = []struct{
	Name string
	Count int
} {
	{"lalala", 1},
	{"bububu", 2},
	{"jujuju", 3},
}

func TestOften001(t * testing.T) {
	u := NewOften(65536)
	for _, a := range data {
		for i := 0; i < a.Count; i++ {
			u.Add(a.Name)
		}
	}
	if u.Count() != 3 {
		t.Fatalf("Count: %v", u.Count())
	}
	
	var temp Results_t
	u.List(&temp, len(data))
	
	if len(temp) != len(data) {
		t.Fatalf("Result: %v", len(temp))
	}
	if temp[0].Value.(string) != data[2].Name {
		t.Fatalf("temp[0]")
	}
	if temp[1].Value.(string) != data[1].Name {
		t.Fatalf("temp[1]")
	}
	if temp[2].Value.(string) != data[0].Name {
		t.Fatalf("temp[2]")
	}
}
