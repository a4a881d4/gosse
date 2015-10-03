package cpbuf

import "testing"
import "unsafe"

func TestSize(t *testing.T) {
	r := ResMem{}
	if unsafe.Sizeof(r) != 1024*1024 {
		t.Error("size error", unsafe.Sizeof(r))
		t.FailNow()
	}
}
