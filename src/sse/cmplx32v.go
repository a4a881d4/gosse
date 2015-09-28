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
import "math"
import "cmplxv"

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
		ret.v = make([]Cmplx32, 0)
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

func (self *Cmplx32) ToComplex() complex128 {
	var r complex128 = complex(float64(self.r), float64(self.i))
	return r
}

func (self *Cmplx32v) ToComplex() []complex128 {
	r := make([]complex128, len(self.v))
	for k, x := range self.v {
		r[k] = x.ToComplex()
	}
	return r
}

func double2Short(y float64) C.short {
	d := math.Floor(0.5 + y)
	if d > 32767.0 {
		return 32767
	}
	if d < -32768.0 {
		return -32768
	}
	return C.short(d)
}

func ToCmplx32(x complex128) Cmplx32 {
	var r Cmplx32
	r.r = double2Short(real(x))
	r.i = double2Short(imag(x))
	return r
}
func ToM128Buf(d []complex128) *Cmplx32v {
	r := NewCmplx32Vec(len(d))
	for k, x := range d {
		r.v[k] = ToCmplx32(x)
	}
	return r
}

func (buf *Cmplx32v) FromBuf() *cmplxv.ComplexV {
	var r cmplxv.ComplexV = cmplxv.ComplexV(buf.ToComplex())
	return &r
}
