package sse

/*
#cgo CFLAGS: -O2 -msse4 -march=core2
#include <xmmintrin.h> // SSE
#include <emmintrin.h> // SSE 2
#include <pmmintrin.h> // SSE 3
#include <tmmintrin.h> // SSSE 3
#include <smmintrin.h> // SSE 4 for media


void xcorr( void *ina, void *inb, int shift, int l, int *out )
{
	int i;
	__m128i *pb;
	__m128i m128_t0,m128_t1,m128_t2,m128_t3,m128_t4,m128_t5,m128_t6,m128_t7,m128_t8,m128_t9,m128_t10,m128_t11,m128_t12;
	__m128i sumi, sumq;
	int *ii,*iq;
	int s0,s1,s2;
	int *pa = (int *)ina;
	pa += shift;
	pb = (__m128i*)inb;

	__m128i  IQ_switch = _mm_setr_epi8(2,3,0,1,6,7,4,5,10,11,8,9,14,15,12,13);
	__m128i  Neg_I = _mm_setr_epi8(0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF);
	__m128i  Neg_R = _mm_setr_epi8(0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1);

	sumi = _mm_xor_si128( sumi, sumi );
	sumq = _mm_xor_si128( sumq, sumq );

	for( i=0;i<l/4;i++ )
	{
		m128_t0 = _mm_loadu_si128((__m128i*)pa); pa+=4;
		m128_t2 = _mm_load_si128(pb);pb++;

		m128_t1 = _mm_shuffle_epi8(m128_t0, IQ_switch);
    	m128_t3 = _mm_sign_epi16  (m128_t2, Neg_I);

    	m128_t8 = _mm_madd_epi16  (m128_t0, m128_t2);
    	m128_t12  = _mm_madd_epi16(m128_t1, m128_t3);

  		sumi = _mm_add_epi32(sumi,m128_t8);
  		sumq = _mm_add_epi32(sumq,m128_t12);
  	}

	ii = (int *)&sumi;
	iq = (int *)&sumq;
	out[0] = 0;
	out[1] = 0;
	for( i=0;i<4;i++ ) {
  		out[0] += ii[i];
  		out[1] += iq[i];
	}
}

*/
import "C"

func (self *Cmplx32v) xcorr(in *Cmplx32v, shift int) complex64 {
	l := len(self.v)
	var ret [2]C.int
	if shift+l > len(in.v) {
		l = len(in.v) - shift
	}
	C.xcorr(in.d, self.d, C.int(shift), C.int(l), &ret[0])
	return complex(float32(ret[0]), float32(ret[1]))
}
