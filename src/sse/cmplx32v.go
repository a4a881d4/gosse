package sse

/*
#include<stdlib.h>
#include"align.h"
*/

import "C"
import "unsafe"
import "reflect"

type cmplx32 struct {
	r C.short
	i C.short
}
type cmplx32v struct {
	v   []cmplx32
	raw unsafe.Pointer
}

func newCmplx32Vec(l int32) *cmplx32v {
	ret := &cmplx32v()
	p, err := C.align_malloc(&ret.raw, l*4, 16)
	runtime.SetFinalizer(ret, func(ret *complx32v) {
		fmt.Println("gc:", ret.raw)
		C.free(ret.raw)
	})
	h := (*reflect.SliceHeader)(unsafe.Pointer(&ret.v))
	h.Cap = l
	h.Len = l
	h.Data = uintptr(p)
	return ret
}
