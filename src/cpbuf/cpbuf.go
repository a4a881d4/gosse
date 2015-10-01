package cpbuf

/*
#include <string.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <unistd.h>
#include "cpbuf.h"
int allocMem( const char *name, int64 mRes, int64 mCP, int64 mSize, void **pRes, void **pCP, void **pStart )
{
  int error = -5;
  void *ptry;
  void *mpRes, *mpCP, *mpStart;
  mpRes = mpCP = mpStart = (void *)-1;
  int fd = open(name,O_RDWR);
  if (fd<0)
  {
    error = -4;
    goto ERRF;
  }
  if( mRes!=0 )
  {
    mpRes = mmap(0,(size_t)mRes, PROT_READ|PROT_WRITE,MAP_SHARED,fd,0);
    if( mpRes==(void *)-1 )
    {
      error = -1;
      goto ERR;
    }
  }
  ptry = mmap(0,(size_t)(mSize+mCP), PROT_READ|PROT_WRITE,MAP_PRIVATE|MAP_ANON,-1,0);
  munmap( ptry, (size_t)(mSize+mCP) );
  if( ptry!=(void *)-1 )
  {
    mpStart = mmap(ptry,(size_t)mSize, PROT_READ|PROT_WRITE,MAP_SHARED,fd,mRes);
    mpCP = mmap(ptry+(size_t)mSize, (size_t)mCP, PROT_READ|PROT_WRITE,MAP_SHARED|MAP_FIXED,fd,mRes);
  }
  else
  {
    error = -2;
    goto ERR;
  }
  if( mpStart == ptry && mpCP == ptry+(size_t)mSize )
  {
    *pStart = mpStart;
    *pCP = mpCP;
    *pRes = mpRes;
    close(fd);
    return 0;
  }
ERR:
  if( mpRes!=(void *)-1 )
    munmap( mpRes, (size_t)mRes );
  if( mpCP!=(void *)-1 )
    munmap( mpCP, (size_t)mCP );
  if( mpStart!=(void *)-1 )
    munmap( mpStart, (size_t)mSize );
  close(fd);
ERRF:
  return error;
}
*/
import "C"
import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"
)

type CPBuffer struct {
	Size, CP, Res     uint64
	pStart, pCP, pRes unsafe.Pointer
	valid             bool
	Name              string
	meta              *C.CPBMeta
}

func getTempDir() string {
	return "/tmp"
}
func (r *CPBuffer) getFileName() string {
	return getTempDir() + "/" + r.Name
}
func NewCPBuffer(size, cp, res uint64, name string) *CPBuffer {
	r := &CPBuffer{valid: false}
	r.Size, r.CP, r.Res = alSize(size), alSize(cp), alSize(res)
	r.Name = name
	fn := r.getFileName()
	if r.checkFile() {
		r.allocMem(fn)
	}
	return r
}
func bufFromName(name string) *CPBuffer {
	r := &CPBuffer{Name: name}
	if r.loadMeta() == 0 {
		r.allocMem(r.getFileName())
		return r
	}
	return nil
}
func (self *CPBuffer) fixFileSize() bool {
	fn := self.getFileName()
	fd, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0666)
	buf := make([]byte, 16)
	l, _ := fd.WriteAt(buf, int64(self.Size+self.Res-16))
	fd.Close()
	if l == 16 {
		return true
	} else {
		return false
	}
}

func (self *CPBuffer) checkFile() bool {
	fn := self.getFileName()
	fst, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return self.fixFileSize()
	}
	if fst.Size() == int64(self.Size+self.Res) {
		return true
	} else {
		return self.fixFileSize()
	}
}

func (self *CPBuffer) Valid() bool {
	return self.valid
}

func alSize(sz uint64) uint64 {
	var pagesize uint64 = uint64(C.getpagesize())
	l := sz + pagesize - 1
	l -= l & (pagesize - 1)
	return l
}
func (r *CPBuffer) reSetMeta() {
	r.meta = (*C.CPBMeta)(r.pRes)
	C.memset(r.pRes, 0, 256)
	r.meta.cpbLen.resLen = C.int64(r.Res)
	r.meta.cpbLen.dataLen = C.int64(r.Size)
	r.meta.cpbLen.cpLen = C.int64(r.CP)
	C.strcpy((*C.char)(unsafe.Pointer(&(r.meta.name[0]))), C.CString(r.Name))
	r.checkMeta(1)
}
func nameEqual(a string, b []byte) bool {
	for k, x := range []byte(a) {
		if b[k] != x {
			return false
		}
	}
	return true
}

func getMetaFromfile(name string) (*C.CPBMeta, error) {
	fn := getTempDir() + "/" + name
	f, _ := os.Open(fn)
	b := make([]byte, 256)
	l, _ := f.Read(b)
	if l != 256 {
		return nil, errors.New("read file failure")
	}
	if !checkMetaMd5(b) {
		return nil, errors.New("meta check failure")
	}
	p := (*C.CPBMeta)(unsafe.Pointer(&(b[0])))
	if !nameEqual(name, b) {
		return nil, errors.New("name no match")
	}
	return p, nil
}

func (r *CPBuffer) loadMeta() int {
	p, err := getMetaFromfile(r.Name)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	r.Res = uint64(p.cpbLen.resLen)
	r.Size = uint64(p.cpbLen.dataLen)
	r.CP = uint64(p.cpbLen.cpLen)
	return 0
}

func (r *CPBuffer) allocMem(fn string) bool {
	if C.allocMem(C.CString(fn), C.int64(r.Res), C.int64(r.CP), C.int64(r.Size), &(r.pRes), &(r.pCP), &(r.pStart)) == 0 {
		r.reSetMeta()
		r.valid = true
	} else {
		r.valid = false
	}
	return r.valid
}

func (self *CPBuffer) checkMeta(mode int) bool {
	meta := self.meta
	b := (*[256]byte)(unsafe.Pointer(meta))
	if mode == 1 {
		c := getMetaMd5((*b)[:])
		for k, x := range c {
			meta.key[k] = C.char(x)
		}
	} else {
		return checkMetaMd5((*b)[:])
	}
	return true
}

func getMetaMd5(b []byte) []byte {
	m := md5.New()
	m.Write(b[:256-16])
	c := m.Sum(nil)
	return c
}

func checkMetaMd5(b []byte) bool {
	md5 := getMetaMd5(b)
	return bytes.Equal(b[256-16:256], md5)
}
func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		files = append(files, fi.Name())
	}
	return files, nil
}
func (self *CPBuffer) HexKey() string {
	p := (*[16]byte)(unsafe.Pointer(&(self.meta.key[0])))
	return hex.EncodeToString((*p)[:])
}
func findByKey(key string) string {
	path := getTempDir()
	files, err := ListDir(path)
	if err != nil {
		return ""
	}
	ckey, err := hex.DecodeString(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, name := range files {
		p, err := getMetaFromfile(name)
		if err == nil {
			var flag bool = true
			for k, x := range ckey {
				if p.key[k] != C.char(x) {
					flag = false
					break
				}
			}
			if flag {
				return name
			}
		}
	}
	return ""
}
func (self *CPBuffer) getBuf(uint64 off, uint64 l) (unsafe.Pointer, error) {
	ioff := off % self.Size
	if ioff+l > self.Size+self.Res {
		return unsafe.Pointer(nil), errors.New("large than CP")
	}
	return self.pStart + ioff, nil
}
