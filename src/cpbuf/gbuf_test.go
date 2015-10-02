package cpbuf

import "testing"
import "unsafe"
import "fmt"

type timing struct {
	a, b, c int
	d, e, f float64
}

func TestAttach(t *testing.T) {
	r := NewCPBuffer(307200, 4096, 4096, "tx.d")
	p := NewGBuf(r)
	var s *timing
	us, slock, sunlock := p.Attach(unsafe.Sizeof(*s))
	s = (*timing)(us)
	sunlock()
	slock()
	s.a = 1
	sunlock()
	fmt.Println(s)
}
