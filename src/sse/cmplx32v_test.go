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
