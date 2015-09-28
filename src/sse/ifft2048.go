package sse

/*
#cgo CFLAGS: -O2 -msse4 -march=core2

#include"ifft2048.h"

static WORD16 double2short(double d)
{
    d = floor(0.5 + d);
    if (d >= 32767.0) return 32767;
    if (d <-32768.0) return -32768;
    return (WORD16)d;
}

void init_ifft2048_twiddle_factor(WORD16 * ifft2048_r2048twiddle, WORD16 * ifft2048_r32twiddle, WORD16* ifft2048_r64twiddle)
{
    WORD32 i,j;
    WORD32 twiddle_temp[2048][2];
    double factor= 32768;


    for (i=0;i<32;i++)
        for (j=0;j<64;j++)
        {
            *(WORD16 *)(ifft2048_r2048twiddle + i*2*2*64 +2*2*j +0) = double2short( factor * sin(2.0 * PI / 2048 * i *j));
            *(WORD16 *)(ifft2048_r2048twiddle + i*2*2*64 +2*2*j +1) = double2short( factor * cos(2.0 * PI / 2048 * i *j));
            *(WORD16 *)(ifft2048_r2048twiddle + i*2*2*64 +2*2*j +2) = double2short( factor * cos(2.0 * PI / 2048 * i *j));
            *(WORD16 *)(ifft2048_r2048twiddle + i*2*2*64 +2*2*j +3) = double2short( -1*factor * sin(2.0 * PI / 2048 * i *j));
        }

        memcpy((WORD8*)(&twiddle_temp[0][0]), (WORD8*)ifft2048_r2048twiddle, 2048*2*sizeof(WORD32));


        for (i=0;i<2048/4;i++)
        {
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 0) =  twiddle_temp[i*4][0];
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 4) =  twiddle_temp[i*4][1];
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 1) =  twiddle_temp[i*4 + 1][0];
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 5) =  twiddle_temp[i*4 + 1][1];
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 2) =  twiddle_temp[i*4 + 2][0];
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 6) =  twiddle_temp[i*4 + 2][1];
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 3) =  twiddle_temp[i*4 + 3][0];
            *( (WORD32 *)(ifft2048_r2048twiddle) + 2*4*i + 7) =  twiddle_temp[i*4 + 3][1];
        }

        for (i=0;i<8;i++)
            for (j=0;j<8;j++)
            {
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +0) = double2short( factor * sin(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +1) = double2short( factor * cos(2.0 * PI / 64 * i *j));

                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +2) = double2short( factor * sin(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +3) = double2short( factor * cos(2.0 * PI / 64 * i *j));

                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +4) = double2short( factor * sin(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +5) = double2short( factor * cos(2.0 * PI / 64 * i *j));

                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +6) = double2short( factor * sin(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +7) = double2short( factor * cos(2.0 * PI / 64 * i *j));

                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +8) = double2short( factor * cos(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +9) = double2short( -1*factor * sin(2.0 * PI / 64 * i *j));

                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +10) = double2short( factor * cos(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +11) = double2short( -1*factor * sin(2.0 * PI / 64 * i *j));

                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +12) = double2short( factor * cos(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +13) = double2short( -1*factor * sin(2.0 * PI / 64 * i *j));

                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +14) = double2short( factor * cos(2.0 * PI / 64 * i *j));
                *(WORD16 *)(ifft2048_r64twiddle + i*2*2*4*8 +2*2*4*j +15) = double2short( -1*factor * sin(2.0 * PI / 64 * i *j));

            }


            for (i=0;i<8;i++)
                for (j=0;j<4;j++)
                {
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +0) = double2short( factor * sin(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +1) = double2short( factor * cos(2.0 * PI / 32 * i *j));

                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +2) = double2short( factor * sin(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +3) = double2short( factor * cos(2.0 * PI / 32 * i *j));

                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +4) = double2short( factor * sin(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +5) = double2short( factor * cos(2.0 * PI / 32 * i *j));

                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +6) = double2short( factor * sin(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +7) = double2short( factor * cos(2.0 * PI / 32 * i *j));


                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +8) = double2short( factor * cos(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +9) = double2short( -1*factor * sin(2.0 * PI / 32 * i *j));

                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +10) = double2short( factor * cos(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +11) = double2short( -1*factor * sin(2.0 * PI / 32 * i *j));

                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +12) = double2short( factor * cos(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +13) = double2short( -1*factor * sin(2.0 * PI / 32 * i *j));

                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +14) = double2short( factor * cos(2.0 * PI / 32 * i *j));
                    *(WORD16 *)(ifft2048_r32twiddle + i*2*2*4*4 +2*2*4*j +15) = double2short( -1*factor * sin(2.0 * PI / 32 * i *j));

                }

}

void ifft2048(__m128i *InBuf, __m128i *OutBuf, __m128i *ifft2048_r32twiddle, __m128i *ifft2048_r64twiddle, __m128i *ifft2048_r2048twiddle)
{
    WORD32 i;
    __m128i m128_t0,m128_t1,m128_t2,m128_t3,m128_t4,m128_t5,m128_t6,m128_t7,m128_t8,m128_t9,m128_t10,m128_t11,m128_t12;
    __m128i OutTmp[512];
    __m128i ifft2048_Temp64_32_Buf[64];
    __m128i TransposeBuf[64];
    WORD32 in_span, out_span;

    __m128i  IQ_switch = _mm_setr_epi8(2,3,0,1,6,7,4,5,10,11,8,9,14,15,12,13);
    __m128i  Neg_I = _mm_setr_epi8(0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF);
    __m128i  Neg_R = _mm_setr_epi8(0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1);
    __m128i  Const_0707 = _mm_setr_epi8(0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A,0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A);
    __m128i  Const_0707_Minus = _mm_setr_epi8(0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5);

    in_span = 16;
    for (i=0;i<16;i++)
    {
        radix32((__m128i *)(InBuf) + i, in_span, (__m128i *)ifft2048_Temp64_32_Buf, (__m128i *)ifft2048_r32twiddle, (__m128i *)(OutTmp) + i);
    }

    DIV_4_2((__m128i *)(OutTmp) + 0);
    DIV_4_2((__m128i *)(OutTmp) + 4);
    DIV_4_2((__m128i *)(OutTmp) + 8);
    DIV_4_2((__m128i *)(OutTmp) + 12);
    for (i=1;i<4;i++)
    {
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 0, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 0) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 1, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 1) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 2, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 2) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 3, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 3) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 4, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 4) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 5, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 5) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 6, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 6) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 7, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 7) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 8, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 8) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 9, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 9) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 10,(__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 10) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 11, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 11) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 12, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 12) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 13, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 13) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 14,(__m128i *) ifft2048_r2048twiddle + 2*(i*16 + 14) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 15, (__m128i *)ifft2048_r2048twiddle + 2*(i*16 + 15) );
    }

    in_span = 1;
    out_span = 8;
    transpose64x4_mulzero((__m128i *)(OutTmp) + 0*64, in_span, (__m128i *)TransposeBuf);
    radix64((__m128i *)TransposeBuf, in_span,(__m128i *)OutBuf + 0, out_span, (__m128i *)ifft2048_r64twiddle);
    for (i=1;i<8;i++)
    {
        transpose64x4((__m128i *)(OutTmp) + i*64, in_span, (__m128i *)TransposeBuf, (__m128i *)ifft2048_r2048twiddle + 2*i*64);
        radix64((__m128i *)TransposeBuf, in_span,(__m128i *)OutBuf + i, out_span, (__m128i *)ifft2048_r64twiddle);
    }
}

*/
import "C"

type IFFT2048Plan struct {
	r2048twiddle *Cmplx32v
	r32twiddle   *Cmplx32v
	r64twiddle   *Cmplx32v
}

func NewIFFT2048Plan() *IFFT2048Plan {
	r := &IFFT2048Plan{}
	r.r2048twiddle = NewCmplx32Vec(2048)
	r.r32twiddle = NewCmplx32Vec(32 * 4)
	r.r64twiddle = NewCmplx32Vec(64 * 4)
	C.init_ifft2048_twiddle_factor((*C.WORD16)(r.r2048twiddle.d), (*C.WORD16)(r.r32twiddle.d), (*C.WORD16)(r.r64twiddle.d))
	return r
}

func (self *IFFT2048Plan) Do(in, out *Cmplx32v) {
	C.ifft2048((*C.__m128i)(in.d), (*C.__m128i)(out.d), (*C.__m128i)(self.r32twiddle.d), (*C.__m128i)(self.r64twiddle.d), (*C.__m128i)(self.r2048twiddle.d))
}
