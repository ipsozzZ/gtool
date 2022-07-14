package sets

import (
	"testing"
)

func TestAA(t *testing.T) {
	var a *II
	var res1, res2 = a.calII()
	t.Logf("res: %v, %v \n", res1, res2)
}

type II struct {
	next  *II
	value int32
}

func (i *II) calII() (v *II, ex bool) {
	if i == nil || i.next == nil {
		//v = -1
		return nil, false
	}
	v = i.next
	ex = true
	return
}
