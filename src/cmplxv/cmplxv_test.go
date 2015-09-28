package cmplxv

import "math"
import "testing"
import "fmt"

func TestComplexV(t *testing.T) {
	a := Zeros(10)
	fmt.Println(a)
	b := Arange(10)
	c := b.Mul(2.5).Cmul(complex(0., math.Pi))
	fmt.Println(c)
	d := c.Exp()
	fmt.Println(d)
}
