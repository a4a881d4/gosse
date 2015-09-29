package sse

import "math"

//import "math/cmplx"
import "testing"
import "cmplxv"

//import "fmt"

func TestFFT2048(t *testing.T) {
	aSin := cmplxv.Arange(2048).Mul(0. / 2048.).Cmul(complex(0., math.Pi*2)).Exp().Mul(64.)
	in := ToM128Buf(*aSin)
	fplan := NewFreqCorrectPlan()
	plan := NewFFT2048Plan()
	iplan := NewIFFT2048Plan()
	fSin := NewCmplx32Vec(2048)
	ifSin := NewCmplx32Vec(2048)

	for i := 0; i < 2048; i++ {
		iplan.Do(in, fSin)
		fplan.DoA(fSin, fSin, 0, 0, 1024)
		plan.Do(fSin, ifSin)
		rSin := ifSin.FromBuf()
		err := rSin.Vsub(in.FromBuf().Mul(1.))
		fplan.Do(in, in, 0, 4*65536)
		avg := math.Sqrt(err.Power() / 2048.)
		if avg > 8. {
			inx, peak := err.FindMax()
			finx, fpeak, _ := fSin.FindMax()
			t.Error(i, finx, fpeak, inx, peak, (*err)[inx], in.v[inx], fSin.v[i], ifSin.v[inx])

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
