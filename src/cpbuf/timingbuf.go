package cpbuf

/*
#include "cpbuf.h"
*/
import "C"
import (
	"cpuclock"
	"reflect"
	"unsafe"
)

type TimingBuf struct {
	mem  *CPBuffer
	clk  []cpuclock.ClockConv
	head *C.TimingBufHead
	key  string
}

func NewTimingBuf(m *CPBuffer) *TimingBuf {
	r := &TimingBuf{mem: m}
	r.head = (*C.TimingBufHead)(m.pRes)
	r.key = m.HexKey()
	h := (*reflect.SliceHeader)(unsafe.Pointer(&(r.clk)))
	h.Cap = 1
	h.Len = 1
	h.Data = uintptr(unsafe.Pointer(&(r.head.buf[0])))
	r.clk[0].Init()
	return r
}

func TimingBufFromFile(name string) *TimingBuf {
	buf := bufFromName(name)
	if buf != nil {
		return NewTimingBuf(buf)
	}
	return nil
}

func TimingBufFromKey(key string) *TimingBuf {
	name := findByKey(key)
	if name != "" {
		return TimingBufFromFile(name)
	}
	return nil
}

func (c *TimingBuf) ClkRoutine() {
	c.clk[0].Stop = false
	cpuclock.Routine(&(c.clk[0]))
}
func (c *TimingBuf) ClkStop() {
	c.clk[0].Stop = true
}
