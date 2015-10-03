package cpbuf

import (
	"fmt"
	"reflect"
	"unsafe"
)

type beShared struct {
	a, b, c int
	d, e, f float64
}

func (r beShared) setA(k int) {
	r.a = k
}

type demoBuf struct {
	GBuf
	my SharedObj
}

func NewDemoBuf(m *CPBuffer) *demoBuf {
	r := &demoBuf{}
	r.New(m)
	return r
}

func (r *demoBuf) Attach() {
	var b *beShared
	r.my = r.GBuf.Attach(unsafe.Sizeof(*b))
	r.my.Obj = (*beShared)(r.my.Obj.(unsafe.Pointer))
}

func (r *demoBuf) Dump() {
	var p *beShared = r.my.Obj.(*beShared)
	fmt.Println(*p)
}

func (r *demoBuf) SetA(l int) {
	//r.my.Obj.(beShared).a = l
}
