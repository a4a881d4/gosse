package cpuclock

import "testing"
import "time"
import "fmt"

func TestCPUClock(t *testing.T) {
	start_cpu := getNow()
	start_sys := time.Now().UnixNano()
	time.Sleep(time.Millisecond * 1000)
	end_cpu := getNow()
	end_sys := time.Now().UnixNano()
	fmt.Println(end_cpu - start_cpu)
	fmt.Println(end_sys - start_sys)
	c := NewClock()
	go routine(c)
	for c.cnt < 400 {

	}
	c.stop = true
}
