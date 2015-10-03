package cpbuf

/*
#include "cpbuf.h"
#include <libkern/OSAtomic.h>
void raw_spin_lock(OSSpinLock *lock){
	OSSpinLockLock(lock);
}
void raw_spin_unlock(OSSpinLock *lock){
	OSSpinLockUnlock(lock);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type GBuf struct {
	mem    *CPBuffer
	head   *C.GBufHead
	key    string
	top    uintptr
	buf    uintptr
	layout map[string]SharedObj
}
type Buf interface {
	New(m *CPBuffer)
	FromFile(name string)
	FromKey(key string)
}

func (r *GBuf) New(m *CPBuffer) {
	r.mem = m
	r.head = (*C.GBufHead)(m.pRes)
	r.key = m.HexKey()
	r.top = uintptr(0)
	r.buf = uintptr(unsafe.Pointer(&(r.head.buf[0])))
}
func (r *GBuf) FromFile(name string) {
	buf := bufFromName(name)
	if buf != nil {
		r.New(buf)
	} else {
		fmt.Println("miss file:", name)
	}
}
func (r *GBuf) FromKey(key string) {
	name := findByKey(key)
	if name != "" {
		r.FromFile(name)
	}
	fmt.Println("miss key")
}

func NewGBuf(m *CPBuffer) *GBuf {
	r := &GBuf{}
	r.New(m)
	return r
}

func (r *GBuf) getPointer(off uintptr) unsafe.Pointer {
	return unsafe.Pointer(r.buf + off)
}

type SharedObj struct {
	Obj    interface{}
	lock   func()
	unlock func()
}

func (r *GBuf) Attach(l uintptr, key string) SharedObj {
	ret := r.getPointer(r.top)
	l += 15
	l -= l & 15
	lock := (*C.OSSpinLock)(r.getPointer(r.top + l))
	flock := func() {
		C.raw_spin_lock(lock)
	}
	funlock := func() {
		C.raw_spin_unlock(lock)
	}
	l = l + 16
	pkey := (*[32]byte)(r.getPointer(r.top + l))
	l += 32
	for k, _ := range pkey {
		pkey[k] = ([]byte(key))[k]
	}
	r.top += l
	retobj := SharedObj{lock: flock, unlock: funlock}
	retobj.Obj = ret
	fmt.Println("obj", retobj)
	return retobj
}
