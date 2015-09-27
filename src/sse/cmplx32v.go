package sse

/*
#include<stdlib.h>
#include<stdio.h>
void *align_malloc( void** praw, int l, int a){
	void *raw = malloc(l+a-1);
	void *r = (void *)(((unsigned long)raw+a-1)&(unsigned long)(-a));
	*praw = raw;
	printf("raw=%p,align=%p\n",raw,r);
	return r;
}
*/
import "C"
import "unsafe"
import "reflect"

type Cmplx32 struct {
	r C.short
	i C.short
}
type Cmplx32v struct {
	v   []Cmplx32
	raw unsafe.Pointer
}

func NewCmplx32Vec(l int32) *Cmplx32v {
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
