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

type cmplx32 struct {
	r C.short
	i C.short
}
type cmplx32v struct {
	v   []cmplx32
	raw unsafe.Pointer
}

func newCmplx32Vec(l int) *cmplx32v {
	ret := &cmplx32v{}
	p := C.align_malloc(&ret.raw, C.int(l*4), 16)
	runtime.SetFinalizer(ret, func(ret *cmplx32v) {
		fmt.Println("gc:", ret.raw)
		C.free(ret.raw)
	})
	h := (*reflect.SliceHeader)(unsafe.Pointer(&ret.v))
	h.Cap = l
	h.Len = l
	h.Data = uintptr(p)
	return ret
}
