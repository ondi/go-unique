//
//
//

package unique

import (
	"testing"

	"gotest.tools/assert"
)

func TestUnique01(t *testing.T) {
	u := NewUnique(65536)
	u.AddUint64(1)
	assert.Assert(t, u.Count() == 1)
}
