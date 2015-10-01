package cpuclock

/*
typedef struct { unsigned long t[2]; } timing;
#define timing_now(x) asm volatile(".byte 15;.byte 49" : "=a"((x)->t[0]),"=d"((x)->t[1]))
long long int __getNow()
{
  timing now;
  timing_now(&now);
  return (long long int)now.t[0]+4294967296LL*(long long int)now.t[1];
}

*/
import "C"
import (
	"fmt"
	"time"
)

func getNow() uint64 {
	n := C.__getNow()
	return uint64(n)
}

func (self *clockConv) getSysNow() float64 {
	n := time.Now().UnixNano() - self.start
	return float64(n)
}

type clockConv struct {
	s     float64
	a     float64
	stop  bool
	cnt   int
	start int64
}

func NewClock() *clockConv {
	r := &clockConv{start: time.Now().UnixNano(), stop: false}
	return r
}
func (self *clockConv) Now() float64 {
	return (self.s + self.a*float64(getNow()))
}
func (c *clockConv) init(n int) {
	cpu := make([]float64, n)
	sys := make([]float64, n)
	for i := 0; i < n; i++ {
		cpu[i] = float64(getNow())
		sys[i] = c.getSysNow()
		time.Sleep(time.Millisecond * 1)
	}
	var y [2]float64
	var x [2][2]float64
	y[0], y[1] = 0., 0.
	x[0][0], x[0][1], x[1][0], x[1][1] = 0., 0., 0., 0.
	for k, a := range sys {
		y[0] += a
		y[1] += a * cpu[k]
	}
	for _, a := range cpu {
		x[0][0] += 1.
		x[1][0] += a
		x[1][1] += a * a
	}
	x[0][1] = x[1][0]
	det := x[0][0]*x[1][1] - x[1][0]*x[0][1]
	c.s = 1. / det * (x[1][1]*y[0] - x[1][0]*y[1])
	c.a = 1. / det * (-x[0][1]*y[0] + x[0][0]*y[1])
}
func routine(c *clockConv) {
	c.init(1000)
	sumerr := 0.
	for c.cnt = 0; !c.stop; c.cnt++ {
		err := c.Now() - c.getSysNow()
		sumerr += err
		c.s -= sumerr*0.2 + err*0.5
		c.a -= err*1e-17 + sumerr*3e-18
		fmt.Println(err, c.a, c.s, sumerr)
		time.Sleep(time.Millisecond * 100)
	}
}
