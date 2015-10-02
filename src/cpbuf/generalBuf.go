package cpbuf

/*
#include "cpbuf.h"
void raw_spin_lock(raw_spinlock_t *lock){
	__raw_spin_lock(lock);
}
void raw_spin_unlock(raw_spinlock_t *lock){
	__raw_spin_unlock(lock);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type GBuf struct {
	mem  *CPBuffer
	head *C.GBufHead
	key  string
	top  uintptr
	buf  uintptr
}

func NewGBuf(m *CPBuffer) *GBuf {
	r := &GBuf{mem: m}
	r.head = (*C.GBufHead)(m.pRes)
	r.key = m.HexKey()
	r.top = uintptr(0)
	r.buf = uintptr(unsafe.Pointer(&(r.head.buf[0])))
	return r
}

func GBufFromFile(name string) *GBuf {
	buf := bufFromName(name)
	if buf != nil {
		return NewGBuf(buf)
	}
	fmt.Println("miss file")
	return nil
}

func GBufFromKey(key string) *GBuf {
	name := findByKey(key)
	if name != "" {
		return GBufFromFile(name)
	}
	fmt.Println("miss key")
	return nil
}
func (r *GBuf) getPointer(off uintptr) unsafe.Pointer {
	return unsafe.Pointer(r.buf + off)
}

type SharedObj struct {
	Obj    interface{}
	lock   func()
	unlock func()
}

func (r *GBuf) Attach(l uintptr) (unsafe.Pointer, func(), func()) {
	ret := r.getPointer(r.top)
	l += 15
	l -= l & 15
	lock := (*C.raw_spinlock_t)(r.getPointer(r.top + l))
	flock := func() {
		C.raw_spin_lock(lock)
	}
	funlock := func() {
		C.raw_spin_unlock(lock)
	}
	l = l + 16
	//pkey := (*[32]byte)(r.getPointer(r.top + l))
	l += 32
	//for k, _ := range pkey {
	//	pkey[k] = ([]byte(key))[k]
	//}
	r.top += l
	return ret, flock, funlock
}
