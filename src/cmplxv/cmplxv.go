package cmplxv

import "math"
import "math/cmplx"
import "math/rand"

type ComplexV []complex128

func Zeros(l int) *ComplexV {
	r := ComplexV(make([]complex128, l))
	for k, _ := range r {
		r[k] = complex(0., 0.)
	}
	return &r
}
func Arange(l int) *ComplexV {
	r := ComplexV(make([]complex128, l))
	for k, _ := range r {
		r[k] = complex(float64(k), 0.)
	}
	return &r
}
func (self *ComplexV) Mul(x float64) *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = a * complex(x, 0.)
	}
	return r
}
func (self *ComplexV) Exp() *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = cmplx.Exp(a)
	}
	return r
}
func (self *ComplexV) Pow(x float64) *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = cmplx.Pow(a, complex(x, 0.))
	}
	return r
}
func (self *ComplexV) Conj() *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = cmplx.Conj(a)
	}
	return r
}
func (self *ComplexV) Cmul(x complex128) *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = a * x
	}
	return r
}
func (self *ComplexV) SwapRI() *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = complex(imag(a), real(a))
	}
	return r
}
func (self *ComplexV) Vadd(b *ComplexV) *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = (*b)[k] + a
	}
	return r
}
func (self *ComplexV) Vsub(b *ComplexV) *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = a - (*b)[k]
	}
	return r
}
func (self *ComplexV) Vmul(b *ComplexV) *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = a * (*b)[k]
	}
	return r
}
func (self *ComplexV) VmulConj(b *ComplexV) *ComplexV {
	r := Zeros(len(*self))
	for k, a := range *self {
		(*r)[k] = a * cmplx.Conj((*b)[k])
	}
	return r
}
func (self *ComplexV) Sum() complex128 {
	r := complex(0., 0.)
	for _, a := range *self {
		r += a
	}
	return r
}
func (self *ComplexV) Dot(b *ComplexV) complex128 {
	return self.VmulConj(b).Sum()
}
func (self *ComplexV) Power() float64 {
	r := self.Dot(self)
	return real(r)
}
func CosSin(l, f, a int) *ComplexV {
	aSin := Arange(l).Mul(float64(f) / float64(l)).Cmul(complex(0., math.Pi*2)).Exp().Mul(float64(a))
	return aSin
}
func (self *ComplexV) FindMax() (inx int, a float64) {
	a = 0.
	inx = 0
	for k, x := range *self {
		p := x * cmplx.Conj(x)
		if real(p) > a {
			a = real(p)
			inx = k
		}
	}
	return inx, a
}
func Noise(l int) *ComplexV {
	r := Zeros(l)
	for k, _ := range *r {
		(*r)[k] = complex(rand.NormFloat64(), rand.NormFloat64())
	}
	return r
}
