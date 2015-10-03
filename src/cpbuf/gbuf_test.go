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
	fmt.Println(p)
	s := &timing{}
	obj := p.Attach(unsafe.Sizeof(*s))
	fmt.Println(obj)
	obj.lock()
	s = (*timing)(obj.Obj.(unsafe.Pointer))
	s.a = 1
	obj.unlock()
	fmt.Println(s)
}
func TestDemo(t *testing.T) {
	r := NewCPBuffer(307200, 4096, 4096, "tx.d")
	one := demoBuf{}
	one.New(r)
	one.Attach()
	another := demoBuf{}
	another.FromFile("tx.d")
	(&another).Attach()
	one.Dump()
	var o *beShared = (another.my.Obj.(*beShared))

	fmt.Println(o.b)
	another.my.lock()
	o.b = 2
	another.my.unlock()
	one.Dump()
}
