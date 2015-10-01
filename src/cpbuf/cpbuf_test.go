package cpbuf

import "testing"
import "fmt"

func TestCP(t *testing.T) {
	r := NewCPBuffer(307200, 4096, 4096, "tx.d")
	key := r.HexKey()
	fmt.Println("key", key)
	name := findByKey(key)
	fmt.Println("name", name)
	b := bufFromName(name)
	fmt.Println("r", r, "b", b)
	t.Log(r)
	t.Log(b)
}
