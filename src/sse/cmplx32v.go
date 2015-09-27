package sse

/*
#include<stdlib.h>
#include"align.h"
*/
import "C"
import "unsafe"
import "reflect"
import "runtime"
import "fmt"

type Cmplx32 struct {
	r C.short
	i C.short
}
type Cmplx32v struct {
	v   []Cmplx32
	raw unsafe.Pointer
	d   unsafe.Pointer
}

func NewCmplx32Vec(l int) *Cmplx32v {
	ret := &Cmplx32v{}
	ret.d = C.align_malloc(&ret.raw, C.int(l*4), 16)
	runtime.SetFinalizer(ret, func(ret *Cmplx32v) {
		fmt.Println("gc:", ret.raw)
		C.free(ret.raw)
	})
	h := (*reflect.SliceHeader)(unsafe.Pointer(&ret.v))
	h.Cap = l
	h.Len = l
	h.Data = uintptr(ret.d)
	return ret
}

func (self *Cmplx32v) Get() unsafe.Pointer {
	return self.d
}

func (self *Cmplx32) ToComplex() complex64 {
	var r complex64 = complex(float32(self.r), float32(self.i))
	return r
}

func (self *Cmplx32v) ToComplex() []complex64 {
	r := make([]complex64, len(self.v))
	for k, x := range self.v {
		r[k] = x.ToComplex()
	}
	return r
}

func ToM128Buf(d []complex64) *Cmplx32v {
	r := NewCmplx32Vec(len(d))
	for k, x := range d {
		r.v[k].r = C.short(real(x))
		r.v[k].i = C.short(imag(x))
	}
	return r
}
