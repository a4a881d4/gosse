package sse

import "math"

//import "math/cmplx"
import "testing"
import "cmplxv"

import "fmt"

func TestFFT2048(t *testing.T) {
	fplan := NewFreqCorrectPlan()
	plan := NewFFT2048Plan()
	iplan := NewIFFT2048Plan()
	fSin := NewCmplx32Vec(2048)
	ifSin := NewCmplx32Vec(2048)

	for i := 0; i < 128; i++ {
		aSin := cmplxv.Noise(2048).Mul(128.)
		in := ToM128Buf(*aSin)
		plan.Do(in, fSin)
		fmt.Println("fSin", fSin.v[:5])
		fplan.DoA(fSin, fSin, 0, 0, 8192)
		fmt.Println("fSin", fSin.v[:5])
		iplan.Do(fSin, ifSin)
		rSin := ifSin.FromBuf()
		x := ifSin.xcorr(in, 0) / in.xcorr(in, 0)
		fmt.Println("x=", x)
		err := rSin.Vsub(in.FromBuf().Cmul(x))
		//fplan.Do(in, in, 0, 4*65536)
		avg := math.Sqrt(err.Power() / 2048.)
		if avg > 4. {
			inx, peak := err.FindMax()
			t.Error(i, inx, peak, (*err)[inx], in.v[inx], ifSin.v[inx], avg)
			t.Error(in.v[:5], ifSin.v[:5])
		}
	}
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
