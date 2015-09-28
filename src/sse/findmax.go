package sse

/*
#cgo CFLAGS: -O2 -msse4 -march=core2
#include <xmmintrin.h> // SSE
#include <emmintrin.h> // SSE 2
#include <pmmintrin.h> // SSE 3
#include <tmmintrin.h> // SSSE 3
#include <smmintrin.h> // SSE 4 for media

void findCmax( void *in, int l, int *iMax, int *avg )
{
	__m128i *ibuf,*obuf;
	__m128i id,pow,max,sum, index, mask, ones, spow;
	int *sp;
	int *ip;
	short *sip;
	int i,j;
	ones = _mm_set_epi32(4,4,4,4);
	mask = _mm_set_epi32(0xfffff800,0xfffff800,0xfffff800,0xfffff800);
	index = _mm_set_epi32(3,2,1,0);

	ibuf=(__m128i *)in;

	max = _mm_xor_si128(max,max);
	sum = _mm_xor_si128(sum,sum);

	for( i=0;i<l/4;i++ )
	{
		id=_mm_load_si128(ibuf);
		pow = _mm_madd_epi16( id, id );
		spow = _mm_srai_epi32(pow, 8);

		pow = _mm_and_si128( pow, mask );
		pow = _mm_or_si128( pow, index );
		max = _mm_max_epi32( max, pow );
		index = _mm_add_epi32( index, ones );

		sum = _mm_add_epi32(sum,spow);

		ibuf++;
		obuf++;
	}
	*iMax=0;
	*avg=0;
	sp = (int *)&max;
	for( i=0;i<4;i++ )
	{
		if( sp[i]>*iMax )
			*iMax=sp[i];
	}
	ip = (int *)&sum;
	for( i=0;i<4;i++ )
	{
		*avg+=(int)ip[i];
	}
}
*/
import "C"

func (self *Cmplx32v) FindMax() (iinx, imax, iavg int) {
	l := len(self.v)
	var max C.int
	var avg C.int
	C.findCmax(self.d, C.int(l), &max, &avg)
	iinx = int(max) & 0x7ff
	imax = int(max) >> 11
	iavg = int(avg)
	return iinx, imax, iavg
}
