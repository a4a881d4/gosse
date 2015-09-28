package sse

import "math"

//import "math/cmplx"
import "testing"
import "cmplxv"
import "fmt"

func TestFFT2048(t *testing.T) {
	plan := NewFFT2048Plan()
	aSin := cmplxv.Arange(2048).Mul(15. / 2048.).Cmul(complex(0., math.Pi*2)).Exp().Mul(127.)
	fSin := NewCmplx32Vec(2048)
	ifSin := NewCmplx32Vec(2048)
	in := ToM128Buf(*aSin)
	plan.Do(in, fSin)
	fmt.Println(fSin.v[:30])
	iplan := NewIFFT2048Plan()
	iplan.Do(fSin, ifSin)
	rSin := ifSin.FromBuf()
	fmt.Println((*rSin)[:5])
	fmt.Println(math.Sqrt(rSin.Vsub(aSin.Mul(4096./128.)).Power() / 2048.))
}

func TestFFTEachFreq(t *testing.T) {
	plan := NewFFT2048Plan()
	for i := 0; i < 2048; i++ {
		aSin := cmplxv.CosSin(2048, i, 63)
		in := ToM128Buf(*aSin)
		out := NewCmplx32Vec(2048)
		plan.Do(in, out)
		inx, max, avg := out.FindMax()
		if inx != i {
			t.Error(i, inx, max, avg)
			t.Fail()
		}
	}
}
