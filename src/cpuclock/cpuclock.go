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

func (c *clockConv) getNow() int64 {
	n := C.__getNow()
	return int64(n) - c.off
}

func (self *clockConv) getSysNow() float64 {
	n := time.Now().UnixNano() - self.start
	return float64(n)
}

type clockConv struct {
	s     float64
	a     float64
	Stop  bool
	Cnt   int
	start int64
	y     float64
	x     float64
	xy    float64
	xx    float64
	alf   float64
	off   int64
	ss    int64
	xs    [2]float64
	cpu_s int64
	sys_s int64
	lunch int64
}

func NewClock() *clockConv {
	r := &clockConv{lunch: time.Now().UnixNano(), Stop: false, alf: 0.5}
	r.y, r.xy = 0., 0.
	r.x, r.xx = 0., 0.
	r.off = 0
	r.ss = 1000000
	r.xs[0], r.xs[1] = 0., 0.
	r.start = r.lunch
	return r
}
func (self *clockConv) Now() float64 {
	return (self.s + self.a*float64(self.getNow()) + float64(self.start))
}
func (c *clockConv) Run() (r, s int64) {
	r = int64(c.s + c.a*float64(c.getNow()))
	r += (c.start - c.lunch)
	return r, c.lunch
}
func (c *clockConv) scale(sys float64) {
	o := sys / c.a
	c.off += int64(o)
	c.start += int64(sys)
}
func (c *clockConv) update() float64 {
	cpu := float64(c.getNow())
	sys := c.getSysNow()
	c.y += c.alf * (sys - c.y)
	c.xy += c.alf * (sys*cpu - c.xy)
	c.x += c.alf * (cpu - c.x)
	c.xx += c.alf * (cpu*cpu - c.xx)
	cpu = c.s + c.a*cpu
	err := cpu - sys
	alf := err * err / 1e6
	if alf < 0.5 {
		c.alf = alf
	} else {
		c.alf = 0.5
	}
	ss := int64(10000000. / c.alf)
	if ss < 1000000000 {
		c.ss = ss
	} else {
		c.ss = 1000000000
	}

	return err
}
func (c *clockConv) init(n int) {
	c.cpu_s = int64(C.__getNow())
	c.sys_s = time.Now().UnixNano()
	for i := 0; i < n; i++ {
		c.update()
		time.Sleep(time.Millisecond * 1)
	}
	c.calc()
}
func (c *clockConv) calc() {
	var x [2][2]float64
	var y [2]float64
	y[0], y[1] = c.y, c.xy
	x[0][0], x[1][0], x[0][1], x[1][1] = 1., c.x, c.x, c.xx
	det := x[0][0]*x[1][1] - x[1][0]*x[0][1]
	s := 1. / det * (x[1][1]*y[0] - x[1][0]*y[1])
	a := 1. / det * (-x[0][1]*y[0] + x[0][0]*y[1])
	c.s = s
	c.a = a
}
func (c *clockConv) loop(e float64) {
	c.xs[0] += e
	c.xs[1] += c.xs[0]
	c.s -= 1. * (e*0.3 + c.xs[0]*0.1 + c.xs[1]*0.01)
	//c.a -= 1.e-16 * (e*5. + c.xs[0]*2. + c.xs[1]*0.5)
}
func (c *clockConv) recalc() {
	cpu_n := int64(C.__getNow())
	sys_n := time.Now().UnixNano()
	//fmt.Println("re calc", sys_n-c.sys_s)
	c.a = float64(sys_n-c.sys_s) / float64(cpu_n-c.cpu_s)
	c.s = 0.
	c.xs[0], c.xs[1] = 0., 0.
	c.off = cpu_n
	c.start = sys_n
	c.cpu_s = cpu_n
	c.sys_s = sys_n
}
func Routine(c *clockConv) {
	c.init(1000)
	fmt.Println("time routine start!!!")
	for c.Cnt = 1; !c.Stop; c.Cnt++ {
		err := c.update()
		c.loop(err)
		if c.Cnt%10 == 0 {
			//fmt.Printf("%8.1f %16.10e %16.10e %16.10e %16.10e\n", err, c.a, c.s, c.xs[0], c.xs[1])
		}
		if c.Cnt%1000 == 0 {
			c.recalc()
		}
		time.Sleep(time.Duration(c.ss))
	}
}
