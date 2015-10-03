package cpbuf

/*
#include "cpbuf.h"
int64 allocMemAlign( VMemHead *head, int64 len, int align )
{
	int64 _align = (int64)align;
	__raw_spin_lock(&(head->lockM));
	int64 last = ((head->_brk+_align-1)/_align)*_align;
	head->_brk=last+len;
	__raw_spin_unlock(&(head->lockM));
	return last;
}
void resetVMem( VMemHead *head )
{
	__raw_spin_unlock(&(head->lockM));
}
*/
import "C"
import (
	"unsafe"
)

type VMem struct {
	mem  *CPBuffer
	head *C.VMemHead
	key  string
}

type GMem struct {
	key    string
	off    uint64
	length uint64
}

func (r *VMem) Start() {
	r.head._brk = 0
	C.resetVMem(r.head)
}

func (r *VMem) ReStart() {
	C.resetVMem(r.head)
}

func NewVMem(m *CPBuffer) *VMem {
	r := &VMem{mem: m}
	r.head = (*C.VMemHead)(m.pRes)
	r.key = m.HexKey()
	return r
}

func VMemFromFile(name string) *VMem {
	buf := bufFromName(name)
	if buf != nil {
		return NewVMem(buf)
	}
	return nil
}
func VMemFromKey(key string) *VMem {
	name := findByKey(key)
	if name != "" {
		return VMemFromFile(name)
	}
	return nil
}
func (r *VMem) alloc(l uint64) (unsafe.Pointer, error) {

	return r.allocAlign(l, 16)
}

func (r *VMem) allocAlign(l uint64, a int) (unsafe.Pointer, error) {
	p := C.allocMemAlign(r.head, C.int64(l), C.int(a))
	return r.mem.getBuf(uint64(p), l)
}

func (r *VMem) allocGMemAlign(l uint64, a int) *GMem {
	gm := &GMem{key: r.key}
	gm.off = uint64(C.allocMemAlign(r.head, C.int64(l), C.int(a)))
	gm.length = l
	return gm
}

func (r *VMem) allocGMem(l uint64) *GMem {
	return r.allocGMemAlign(l, 16)
}
