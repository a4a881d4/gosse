package sse

/*
#cgo CFLAGS: -O2 -msse4 -march=core2
#include <xmmintrin.h> // SSE
#include <emmintrin.h> // SSE 2
#include <pmmintrin.h> // SSE 3
#include <tmmintrin.h> // SSSE 3
#include <smmintrin.h> // SSE 4 for media


int freqCorrect( void *in, void *out, void *stwiddle_2048, int nco, int freq, int len, short ampl, short amph )
{
	__m128i *pin, *pout;
	__m128i m128_t0,m128_t1,m128_t2,m128_t3,m128_t4,m128_t5,m128_t6,m128_t7,m128_t8,m128_t9,m128_t10,m128_t11,m128_t12;
	__m128i m128_ampl,m128_amph;
	int i,phase,t;
	int *twiddle = (int *)stwiddle_2048;
    __m128i  IQ_switch = _mm_setr_epi8(2,3,0,1,6,7,4,5,10,11,8,9,14,15,12,13);
	__m128i  Neg_I = _mm_setr_epi8(0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF);
	__m128i  Neg_R = _mm_setr_epi8(0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1);

	pin = (__m128i *)in;
	pout = (__m128i *)out;
	m128_ampl = _mm_set_epi16( ampl, ampl, ampl, ampl, ampl, ampl, ampl, ampl );
	m128_amph = _mm_set_epi16( amph, amph, amph, amph, amph, amph, amph, amph );

	for( i=0;i<(len+3)/4;i++ )
	{
		phase = (nco>>16)&0x7ff;
		m128_t0 = _mm_load_si128(pin);  pin = pin + 1;
		t = twiddle[phase];
		m128_t2 = _mm_set_epi32( t, t, t, t );

		m128_t1 = _mm_shuffle_epi8(m128_t0, IQ_switch);
    	m128_t3 = _mm_sign_epi16  (m128_t2, Neg_I);

    	m128_t8 = _mm_madd_epi16  (m128_t0, m128_t2);
    	m128_t12  = _mm_madd_epi16(m128_t1, m128_t3);

    	m128_t8  = _mm_srli_si128 (m128_t8, 2);
    	m128_t12 = _mm_blend_epi16(m128_t12,m128_t8, 0x55);
	    m128_t12 = _mm_mullo_epi16( m128_t12, m128_ampl );
		m128_t12 = _mm_mulhi_epi16( m128_t12, m128_amph );

    	_mm_store_si128(pout, m128_t12);  pout = pout + 1;
    	nco += freq;
	}
	return nco;
}

*/
import "C"
import "cmplxv"

type FreqCorrectPlan struct {
	r2048twiddle *Cmplx32v
}

func NewFreqCorrectPlan() *FreqCorrectPlan {
	t := cmplxv.CosSin(2048, 1, 32767)
	r := &FreqCorrectPlan{ToM128Buf(*t)}
	return r
}

func (self *FreqCorrectPlan) Do(in, out *Cmplx32v, nco, f int) int {
	r := C.freqCorrect(in.d, out.d, self.r2048twiddle.d, C.int(nco), C.int(f), C.int(len(in.v)), C.short(-4), C.short(-32768))
	return int(r)
}
func (self *FreqCorrectPlan) DoA(in, out *Cmplx32v, nco, f, a int) int {
	r := C.freqCorrect(in.d, out.d, self.r2048twiddle.d, C.int(nco), C.int(f), C.int(len(in.v)), C.short(-4), C.short(-a))
	return int(r)
}
