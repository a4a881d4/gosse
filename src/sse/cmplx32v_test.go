package sse

import "testing"
import "fmt"
import "unsafe"

func TestNewCmplx32Vec(t *testing.T) {
	tut := NewCmplx32Vec(128)
	fmt.Println(tut.v, len(tut.v), tut.raw, tut.d)
	tut.v[0].r = 1
	fmt.Println(unsafe.Pointer(&tut.v[127]))
	fmt.Println(tut.Get())
}

func TestToBuf(t *testing.T) {
	var a complex64 = 1. - 1.i
	var z complex64 = 0.
	data := make([]complex64, 128)
	for k, _ := range data {
		data[k] = z
		z = z + a
	}
	buf := ToM128Buf(data)
	fmt.Println(buf.v)
	r := buf.ToComplex()
	fmt.Println(r)
}
