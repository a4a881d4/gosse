package cpbuf

import "unsafe"
import "fmt"

type beShared struct {
	a, b, c int
	d, e, f float64
}

func (r beShared) setA(k int) {
	r.a = k
}

type demoBuf struct {
	buf *GBuf
	my  SharedObj
}

func NewDemoBuf(m *CPBuffer) *demoBuf {
	r := &demoBuf{}
	r.buf = NewGBuf(m)
	return r
}

func DemoBufFromFile(name string) *demoBuf {
	r := &demoBuf{}
	r.buf = GBufFromFile(name)
	return r
}

func (r *demoBuf) Attach() {
	var us unsafe.Pointer
	var b *beShared
	us, r.my.lock, r.my.unlock = r.buf.Attach(unsafe.Sizeof(*b))
	r.my.Obj = *((*beShared)(us))
}

func (r *demoBuf) Dump() {
	fmt.Println(r.my.Obj.(timing))
}

func (r *demoBuf) SetA(l int) {
	//r.my.Obj.(timing).a = l
}
