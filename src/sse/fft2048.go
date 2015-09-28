package sse

/*
#cgo CFLAGS: -O2 -msse4 -march=core2

#include <stdio.h>
#include <stdlib.h>
#include <math.h>

#include <xmmintrin.h> // SSE
#include <emmintrin.h> // SSE 2
#include <pmmintrin.h> // SSE 3
#include <tmmintrin.h> // SSSE 3
#include <smmintrin.h> // SSE 4 for media
#ifndef PI
#define PI (3.14159265358979323846)
#endif

typedef short WORD16;
typedef int WORD32;

#define MUL_CROSS(in_addr, twiddle_addr) \
{ \
    __m128i * twiddle_addr_temp; \
    __m128i * in_addr_temp; \
     \
    in_addr_temp = in_addr; \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
     \
    twiddle_addr_temp = twiddle_addr; \
    m128_t3 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
    m128_t4 = _mm_shuffle_epi8(m128_t3, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
     \
    m128_t1 = _mm_madd_epi16(m128_t0, m128_t3); \
    m128_t0 = _mm_madd_epi16(m128_t0, m128_t4); \
    m128_t0 = _mm_srli_si128(m128_t0, 2); \
    m128_t1 = _mm_blend_epi16(m128_t1,m128_t0, 0x55); \
    _mm_store_si128((__m128i *)in_addr_temp, m128_t1); \
}

#define DIV_4_2_span(in_addr, span); \
{ \
    __m128i * in_addr_temp; \
    in_addr_temp = in_addr; \
     \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
     \
    in_addr_temp = in_addr_temp + span; \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
     \
    in_addr_temp = in_addr_temp + span; \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
     \
    in_addr_temp = in_addr_temp + span; \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
}

#define DIV_4_2(in_addr); \
{ \
    __m128i * in_addr_temp; \
    in_addr_temp = in_addr; \
     \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
     \
    in_addr_temp = in_addr_temp + 1; \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
     \
    in_addr_temp = in_addr_temp + 1; \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
     \
    in_addr_temp = in_addr_temp + 1; \
    m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
    _mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0, 1)); \
}

#define radix4_register(in0, in1, in2, in3, out_addr, out_span) \
{ \
    __m128i * in_addr_temp; \
    m128_t2 = _mm_adds_epi16(in0, in2); \
    m128_t6 = _mm_adds_epi16(in1, in3); \
    m128_t7=  _mm_adds_epi16(m128_t2, m128_t6); \
    m128_t2 = _mm_subs_epi16(m128_t2, m128_t6); \
     \
    m128_t6 =  m128_t7; \
     \
    m128_t3 = _mm_subs_epi16(in0, in2); \
    m128_t12 = _mm_subs_epi16(in1, in3); \
    m128_t12 = _mm_shuffle_epi8(m128_t12, IQ_switch); \
    m128_t12 = _mm_sign_epi16(m128_t12, Neg_R); \
     \
    m128_t7 = _mm_subs_epi16(m128_t3, m128_t12); \
    m128_t3 = _mm_adds_epi16(m128_t3, m128_t12); \
     \
    _mm_store_si128((__m128i *)out_addr, m128_t6); \
    out_addr = out_addr + out_span; \
     \
    _mm_store_si128((__m128i *)out_addr, m128_t7); \
    out_addr = out_addr + out_span; \
     \
    _mm_store_si128((__m128i *)out_addr, m128_t2); \
    out_addr = out_addr + out_span; \
     \
    _mm_store_si128((__m128i *)out_addr, m128_t3); \
}

#define radix8_0(in_addr, in_span, out_addr, out_span) \
{ \
    __m128i * out_addr_temp1, * out_addr_temp2; \
    __m128i * in_addr_temp; \
    in_addr_temp = in_addr; \
    m128_t0 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t1 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t2 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t3 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t4 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t5 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t6 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t7 = _mm_load_si128((__m128i *)in_addr_temp); \
     \
    m128_t8 = _mm_adds_epi16(m128_t0, m128_t4); \
    m128_t9 = _mm_subs_epi16(m128_t0, m128_t4); \
    m128_t10 = _mm_adds_epi16(m128_t1, m128_t5); \
    m128_t11 = _mm_subs_epi16(m128_t1, m128_t5); \
    m128_t0 = _mm_adds_epi16(m128_t2, m128_t6); \
    m128_t1 = _mm_subs_epi16(m128_t2, m128_t6); \
    m128_t4 = _mm_adds_epi16(m128_t3, m128_t7); \
    m128_t5 = _mm_subs_epi16(m128_t3, m128_t7); \
     \
    m128_t1 = _mm_shuffle_epi8(m128_t1, IQ_switch); \
    m128_t1 = _mm_sign_epi16(m128_t1, Neg_I); \
     \
    m128_t2 = _mm_mulhrs_epi16(m128_t11, Const_0707); \
    m128_t3 = _mm_mulhrs_epi16(m128_t5, Const_0707_Minus); \
     \
    m128_t6 = _mm_shuffle_epi8(m128_t2, IQ_switch); \
    m128_t7 = _mm_adds_epi16(m128_t2, m128_t6); \
    m128_t2 = _mm_subs_epi16(m128_t2, m128_t6); \
    m128_t11 = _mm_blend_epi16(m128_t2, m128_t7, 0x55); \
     \
    m128_t6 = _mm_shuffle_epi8(m128_t3, IQ_switch); 	m128_t7 = _mm_adds_epi16(m128_t3, m128_t6); \
    m128_t3 = _mm_subs_epi16(m128_t3, m128_t6); \
    m128_t5 = _mm_blend_epi16(m128_t3, m128_t7, 0xAA); \
     \
     \
    out_addr_temp1 = out_addr; \
    out_addr_temp2 = out_addr + out_span; \
    radix4_register(m128_t8, m128_t10, m128_t0, m128_t4, out_addr_temp1, 2*out_span); \
    radix4_register(m128_t9, m128_t11, m128_t1, m128_t5, out_addr_temp2, 2*out_span); \
}

#define radix4_register_mul_zero(in0, in1, in2, in3, out_addr, out_span, twiddle_addr) \
{ \
    __m128i *in_addr_temp; \
    __m128i *twiddle_addr_temp; \
    m128_t2 = _mm_adds_epi16(in0,in2);  \
    m128_t6 = _mm_adds_epi16(in1,in3);  \
    m128_t7=  _mm_adds_epi16(m128_t2,m128_t6);   \
    m128_t2 = _mm_subs_epi16(m128_t2,m128_t6);   \
    m128_t6 =  m128_t7;  \
     \
    m128_t3 = _mm_subs_epi16(in0,in2);  \
    m128_t12 = _mm_subs_epi16(in1,in3);  \
    m128_t12 = _mm_shuffle_epi8(m128_t12, IQ_switch); \
    m128_t12 = _mm_sign_epi16(m128_t12, Neg_R); \
     \
    m128_t7 = _mm_subs_epi16(m128_t3,m128_t12);  \
    m128_t3 = _mm_adds_epi16(m128_t3,m128_t12);  \
     \
     \
     \
    _mm_store_si128((__m128i *)out_addr,  _mm_srai_epi16(m128_t6,1));   out_addr = out_addr + out_span; \
     \
    twiddle_addr_temp = twiddle_addr + 2; \
    m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 2; \
    m128_t4 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
    m128_t0 = _mm_madd_epi16(m128_t7, m128_t0); \
    m128_t4 = _mm_madd_epi16(m128_t7, m128_t4); \
    m128_t4 = _mm_srli_si128(m128_t4,2); \
    m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
    _mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
     \
    m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 2; \
    m128_t4 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
    m128_t0 = _mm_madd_epi16(m128_t2, m128_t0); \
    m128_t4 = _mm_madd_epi16(m128_t2, m128_t4); \
    m128_t4 = _mm_srli_si128(m128_t4,2); \
    m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
    _mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
     \
    m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
    m128_t4 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
    m128_t0 = _mm_madd_epi16(m128_t3, m128_t0); \
    m128_t4 = _mm_madd_epi16(m128_t3, m128_t4); \
    m128_t4 = _mm_srli_si128(m128_t4,2); \
    m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
    _mm_store_si128((__m128i *)out_addr, m128_t0); \
}

#define radix4_register_mul(in0, in1, in2, in3, out_addr, out_span, twiddle_addr) \
{ \
    __m128i *in_addr_temp; \
    __m128i *twiddle_addr_temp; \
    m128_t2 = _mm_adds_epi16(in0,in2);  \
    m128_t6 = _mm_adds_epi16(in1,in3);  \
    m128_t7=  _mm_adds_epi16(m128_t2,m128_t6);   \
    m128_t2 = _mm_subs_epi16(m128_t2,m128_t6);   \
    m128_t6 =  m128_t7;  \
     \
    m128_t3 = _mm_subs_epi16(in0,in2);  \
    m128_t12 = _mm_subs_epi16(in1,in3); \
    m128_t12 = _mm_shuffle_epi8(m128_t12, IQ_switch); \
    m128_t12 = _mm_sign_epi16(m128_t12, Neg_R); \
     \
    m128_t7 = _mm_subs_epi16(m128_t3,m128_t12);  \
    m128_t3 = _mm_adds_epi16(m128_t3,m128_t12);  \
     \
     \
    twiddle_addr_temp = twiddle_addr; \
    m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 2; \
    m128_t4 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
    m128_t0 = _mm_madd_epi16(m128_t6, m128_t0); \
    m128_t4 = _mm_madd_epi16(m128_t6, m128_t4); \
    m128_t4 = _mm_srli_si128(m128_t4,2); \
    m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
    _mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
     \
    m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 2; \
    m128_t4 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
    m128_t0 = _mm_madd_epi16(m128_t7, m128_t0); \
    m128_t4 = _mm_madd_epi16(m128_t7, m128_t4); \
    m128_t4 = _mm_srli_si128(m128_t4,2); \
    m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
    _mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
     \
    m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 2; \
    m128_t4 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
    m128_t0 = _mm_madd_epi16(m128_t2, m128_t0); \
    m128_t4 = _mm_madd_epi16(m128_t2, m128_t4); \
    m128_t4 = _mm_srli_si128(m128_t4,2); \
    m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
    _mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
     \
    m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
    m128_t4 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
    m128_t4 = _mm_sign_epi16(m128_t4, Neg_I); \
    m128_t0 = _mm_madd_epi16(m128_t3, m128_t0); \
    m128_t4 = _mm_madd_epi16(m128_t3, m128_t4); \
    m128_t4 = _mm_srli_si128(m128_t4,2); \
    m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
    _mm_store_si128((__m128i *)out_addr, m128_t0); \
}

#define radix8_0_mul(in_addr, in_span, out_addr, out_span, twiddle_addr) \
{ \
    __m128i * out_addr_temp1,*out_addr_temp2; \
    __m128i * in_addr_temp; \
    in_addr_temp = in_addr; \
    m128_t0 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t1 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t2 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t3 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t4 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t5 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t6 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t7 = _mm_load_si128((__m128i *)in_addr_temp); \
     \
    m128_t8 = _mm_adds_epi16(m128_t0, m128_t4);  \
    m128_t9 = _mm_subs_epi16(m128_t0, m128_t4);  \
    m128_t10 = _mm_adds_epi16(m128_t1, m128_t5); \
    m128_t11 = _mm_subs_epi16(m128_t1, m128_t5); \
    m128_t0 = _mm_adds_epi16(m128_t2, m128_t6);  \
    m128_t1 = _mm_subs_epi16(m128_t2, m128_t6);  \
    m128_t4 = _mm_adds_epi16(m128_t3, m128_t7);  \
    m128_t5 = _mm_subs_epi16(m128_t3, m128_t7);  \
     \
    m128_t1 = _mm_shuffle_epi8(m128_t1, IQ_switch); \
    m128_t1 = _mm_sign_epi16(m128_t1, Neg_I);  \
     \
    m128_t2 = _mm_mulhrs_epi16(m128_t11, Const_0707);  \
    m128_t3 = _mm_mulhrs_epi16(m128_t5, Const_0707_Minus);   \
     \
    m128_t6 = _mm_shuffle_epi8(m128_t2, IQ_switch);  \
    m128_t7 = _mm_adds_epi16(m128_t2,m128_t6); \
    m128_t2 = _mm_subs_epi16(m128_t2,m128_t6); \
    m128_t11 = _mm_blend_epi16(m128_t2,m128_t7, 0x55); \
     \
    m128_t6 = _mm_shuffle_epi8(m128_t3, IQ_switch);  \
    m128_t7 = _mm_adds_epi16(m128_t3,m128_t6); \
    m128_t3 = _mm_subs_epi16(m128_t3,m128_t6); \
    m128_t5 = _mm_blend_epi16(m128_t3,m128_t7, 0xAA); \
     \
     \
    out_addr_temp1 = out_addr; \
    out_addr_temp2 = out_addr + out_span; \
    radix4_register_mul_zero(m128_t8,m128_t10,m128_t0,m128_t4, out_addr_temp1, 2*out_span,twiddle_addr); \
    radix4_register_mul(m128_t9,m128_t11,m128_t1,m128_t5, out_addr_temp2, 2*out_span,twiddle_addr + 1); \
}

#define radix64(InBuf, in_span, OutBuf, out_span, r64twiddle) \
{ \
    __m128i tmp_buf[64]; \
    radix8_0(InBuf + 0*in_span, 8*in_span, tmp_buf +0*1,8); \
    DIV_4_2_span((__m128i *)(tmp_buf) + 0*8,8); \
    DIV_4_2_span((__m128i *)(tmp_buf) + 4*8,8); \
    radix8_0_mul(InBuf + 1*in_span, 8*in_span, tmp_buf +1, 8, r64twiddle+1*8); \
    radix8_0_mul(InBuf + 2*in_span, 8*in_span, tmp_buf +2, 8, r64twiddle+2*8); \
    radix8_0_mul(InBuf + 3*in_span, 8*in_span, tmp_buf +3, 8, r64twiddle+3*8); \
    radix8_0_mul(InBuf + 4*in_span, 8*in_span, tmp_buf +4, 8, r64twiddle+4*8); \
    radix8_0_mul(InBuf + 5*in_span, 8*in_span, tmp_buf +5, 8, r64twiddle+5*8); \
    radix8_0_mul(InBuf + 6*in_span, 8*in_span, tmp_buf +6, 8, r64twiddle+6*8); \
    radix8_0_mul(InBuf + 7*in_span, 8*in_span, tmp_buf +7, 8, r64twiddle+7*8); \
     \
    radix8_0(tmp_buf + 0*8, 1, OutBuf + 0*out_span, 8*out_span); \
    radix8_0(tmp_buf + 1*8, 1, OutBuf + 1*out_span, 8*out_span); \
    radix8_0(tmp_buf + 2*8, 1, OutBuf + 2*out_span, 8*out_span); \
    radix8_0(tmp_buf + 3*8, 1, OutBuf + 3*out_span, 8*out_span); \
    radix8_0(tmp_buf + 4*8, 1, OutBuf + 4*out_span, 8*out_span); \
    radix8_0(tmp_buf + 5*8, 1, OutBuf + 5*out_span, 8*out_span); \
    radix8_0(tmp_buf + 6*8, 1, OutBuf + 6*out_span, 8*out_span); \
    radix8_0(tmp_buf + 7*8, 1, OutBuf + 7*out_span, 8*out_span); \
}

#define radix4_0_zeromul(in_addr, in_span, out_addr) \
{ \
    __m128i * in_addr_temp; \
    __m128i * out_addr_temp; \
    __m128i * twiddle_addr_temp; \
    in_addr_temp = in_addr; \
    m128_t0 = _mm_load_si128(in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t1 = _mm_load_si128(in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t2 = _mm_load_si128(in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t3 = _mm_load_si128(in_addr_temp); \
    m128_t4 = _mm_adds_epi16(m128_t0,m128_t2);  \
    m128_t5 = _mm_adds_epi16(m128_t1,m128_t3);  \
    m128_t6 = _mm_subs_epi16(m128_t0,m128_t2);  \
    m128_t7 = _mm_subs_epi16(m128_t1,m128_t3);  \
    m128_t7 = _mm_shuffle_epi8(m128_t7, IQ_switch); \
    m128_t7 = _mm_sign_epi16(m128_t7, Neg_R); \
     \
    m128_t8 = _mm_adds_epi16(m128_t4,m128_t5); \
    m128_t8 = _mm_srai_epi16(m128_t8,1); \
    m128_t9 = _mm_subs_epi16(m128_t6,m128_t7); \
    m128_t9 = _mm_srai_epi16(m128_t9,1); \
    m128_t10 = _mm_subs_epi16(m128_t4,m128_t5); \
    m128_t10 = _mm_srai_epi16(m128_t10,1); \
    m128_t11 = _mm_adds_epi16(m128_t6,m128_t7); \
    m128_t11 = _mm_srai_epi16(m128_t11,1); \
     \
    out_addr_temp = out_addr; \
    _mm_store_si128((__m128i *)out_addr_temp, m128_t8);  out_addr_temp = out_addr_temp +1; \
    _mm_store_si128((__m128i *)out_addr_temp, m128_t9);  out_addr_temp = out_addr_temp +1; \
    _mm_store_si128((__m128i *)out_addr_temp, m128_t10); out_addr_temp = out_addr_temp +1; \
    _mm_store_si128((__m128i *)out_addr_temp, m128_t11); \
}

#define radix4_0_mul(in_addr, in_span, out_addr, twiddle_addr) \
{ \
    __m128i * in_addr_temp; \
    __m128i * out_addr_temp; \
    __m128i * twiddle_addr_temp; \
    in_addr_temp = in_addr; \
    m128_t0 = _mm_load_si128(in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t1 = _mm_load_si128(in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t2 = _mm_load_si128(in_addr_temp);  in_addr_temp = in_addr_temp + in_span; \
    m128_t3 = _mm_load_si128(in_addr_temp); \
    m128_t4 = _mm_adds_epi16(m128_t0,m128_t2);  \
    m128_t5 = _mm_adds_epi16(m128_t1,m128_t3);  \
    m128_t6 = _mm_subs_epi16(m128_t0,m128_t2);  \
    m128_t7 = _mm_subs_epi16(m128_t1,m128_t3);  \
    m128_t7 = _mm_shuffle_epi8(m128_t7, IQ_switch); \
    m128_t7 = _mm_sign_epi16(m128_t7, Neg_R);  \
     \
    m128_t8 = _mm_adds_epi16(m128_t4,m128_t5); \
    m128_t9 = _mm_subs_epi16(m128_t6,m128_t7); \
    m128_t10 = _mm_subs_epi16(m128_t4,m128_t5); \
    m128_t11 = _mm_adds_epi16(m128_t6,m128_t7); \
    twiddle_addr_temp = twiddle_addr + 1; \
    m128_t2 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
    m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
    m128_t6 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
    m128_t3 = _mm_shuffle_epi8(m128_t2, IQ_switch); \
    m128_t3 = _mm_sign_epi16(m128_t3, Neg_I); \
    m128_t5 = _mm_shuffle_epi8(m128_t4, IQ_switch); \
    m128_t5 = _mm_sign_epi16(m128_t5, Neg_I); \
    m128_t7 = _mm_shuffle_epi8(m128_t6, IQ_switch); \
    m128_t7 = _mm_sign_epi16(m128_t7, Neg_I); \
     \
    out_addr_temp = out_addr; \
    _mm_store_si128((__m128i *)out_addr_temp, _mm_srai_epi16(m128_t8,1)); out_addr_temp = out_addr_temp +1; \
     \
    m128_t12 = _mm_madd_epi16(m128_t9, m128_t2); \
    m128_t8 = _mm_madd_epi16(m128_t9, m128_t3); \
    m128_t8 = _mm_srli_si128(m128_t8,2); \
    m128_t12 = _mm_blend_epi16(m128_t12,m128_t8, 0x55); \
    _mm_store_si128((__m128i *)out_addr_temp, m128_t12);  out_addr_temp = out_addr_temp +1; \
    m128_t12 = _mm_madd_epi16(m128_t10, m128_t4); \
    m128_t8 = _mm_madd_epi16(m128_t10, m128_t5); \
    m128_t8 = _mm_srli_si128(m128_t8,2); \
    m128_t12 = _mm_blend_epi16(m128_t12,m128_t8, 0x55); \
    _mm_store_si128((__m128i *)out_addr_temp, m128_t12);  out_addr_temp = out_addr_temp +1; \
    m128_t12 = _mm_madd_epi16(m128_t11, m128_t6); \
    m128_t8 = _mm_madd_epi16(m128_t11, m128_t7); \
    m128_t8 = _mm_srli_si128(m128_t8,2); \
    m128_t12 = _mm_blend_epi16(m128_t12,m128_t8, 0x55); \
    _mm_store_si128((__m128i *)out_addr_temp, m128_t12); \
}


#define radix32(InBuf, in_span, Temp64_32_Buf, r32twiddle, OutBuf) \
{ \
    radix4_0_zeromul(InBuf + 0*in_span, 8*in_span, Temp64_32_Buf + 0*4); \
    radix4_0_mul(InBuf + 1*in_span, 8*in_span, Temp64_32_Buf + 1*4, r32twiddle + 4*1); \
    radix4_0_mul(InBuf + 2*in_span, 8*in_span, Temp64_32_Buf + 2*4, r32twiddle + 4*2); \
    radix4_0_mul(InBuf + 3*in_span, 8*in_span, Temp64_32_Buf + 3*4, r32twiddle + 4*3); \
    radix4_0_mul(InBuf + 4*in_span, 8*in_span, Temp64_32_Buf + 4*4, r32twiddle + 4*4); \
    radix4_0_mul(InBuf + 5*in_span, 8*in_span, Temp64_32_Buf + 5*4, r32twiddle + 4*5); \
    radix4_0_mul(InBuf + 6*in_span, 8*in_span, Temp64_32_Buf + 6*4, r32twiddle + 4*6); \
    radix4_0_mul(InBuf + 7*in_span, 8*in_span, Temp64_32_Buf + 7*4, r32twiddle + 4*7); \
     \
    radix8_0(Temp64_32_Buf + 0, 4, OutBuf + 0*in_span, 4*in_span); \
    radix8_0(Temp64_32_Buf + 1, 4, OutBuf + 1*in_span, 4*in_span); \
    radix8_0(Temp64_32_Buf + 2, 4, OutBuf + 2*in_span, 4*in_span); \
    radix8_0(Temp64_32_Buf + 3, 4, OutBuf + 3*in_span, 4*in_span); \
}

#define _MM_TRANSPOSE4_EPI32(in0, in1, in2, in3) \
{ \
    __m128i tmp0, tmp1, tmp2, tmp3; \
    tmp0 =  _mm_unpacklo_epi32(in0, in1); \
    tmp1 =  _mm_unpackhi_epi32(in0, in1); \
    tmp2 =  _mm_unpacklo_epi32(in2, in3); \
    tmp3 =  _mm_unpackhi_epi32(in2, in3); \
 \
    in0 =  _mm_unpacklo_epi64(tmp0, tmp2); \
    in1 =  _mm_unpackhi_epi64(tmp0, tmp2); \
    in2 =  _mm_unpacklo_epi64(tmp1, tmp3); \
    in3 =  _mm_unpackhi_epi64(tmp1, tmp3); \
}

#define transpose64x4_mulzero(InBuf, in_span, OutBuf, twiddle_addr) \
{ \
    __m128i *twiddle_addr_temp; \
    __m128i *in_addr_temp; \
    __m128i *out_addr_temp; \
     \
    for (WORD32 ii=0;ii<64/4;ii++) \
    { \
        in_addr_temp = InBuf + ii*in_span; \
        m128_t1 = _mm_load_si128((__m128i *)in_addr_temp); in_addr_temp = in_addr_temp + 16*in_span; \
        m128_t3 = _mm_load_si128((__m128i *)in_addr_temp); in_addr_temp = in_addr_temp + 16*in_span; \
        m128_t5 = _mm_load_si128((__m128i *)in_addr_temp); in_addr_temp = in_addr_temp + 16*in_span; \
        m128_t7 = _mm_load_si128((__m128i *)in_addr_temp); \
         \
        _MM_TRANSPOSE4_EPI32(m128_t1, m128_t3, m128_t5, m128_t7); \
         \
        out_addr_temp = OutBuf + ii*4; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t1);  out_addr_temp = out_addr_temp + 1; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t3);  out_addr_temp = out_addr_temp + 1; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t5);  out_addr_temp = out_addr_temp + 1; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t7); \
    } \
}

#define transpose64x4(InBuf, in_span, OutBuf, twiddle_addr) \
{ \
    __m128i *twiddle_addr_temp; \
    __m128i *in_addr_temp; \
    __m128i *out_addr_temp; \
     \
    for (WORD32 ii=0;ii<64/4;ii++) \
    { \
        twiddle_addr_temp = twiddle_addr + ii; \
        m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 16; \
        m128_t2 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 16; \
        m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 16; \
        m128_t6 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
        m128_t1 = _mm_shuffle_epi8(m128_t0, IQ_switch); \
        m128_t1 = _mm_sign_epi16(m128_t1, Neg_I); \
        m128_t3 = _mm_shuffle_epi8(m128_t2, IQ_switch); \
        m128_t3 = _mm_sign_epi16(m128_t3, Neg_I); \
        m128_t5 = _mm_shuffle_epi8(m128_t4, IQ_switch); \
        m128_t5 = _mm_sign_epi16(m128_t5, Neg_I); \
        m128_t7 = _mm_shuffle_epi8(m128_t6, IQ_switch); \
        m128_t7 = _mm_sign_epi16(m128_t7, Neg_I); \
         \
         \
        in_addr_temp = InBuf + ii*in_span; \
        m128_t8 = _mm_load_si128((__m128i *)in_addr_temp); in_addr_temp = in_addr_temp + 16*in_span; \
        m128_t9 = _mm_load_si128((__m128i *)in_addr_temp); in_addr_temp = in_addr_temp + 16*in_span; \
        m128_t10 = _mm_load_si128((__m128i *)in_addr_temp); in_addr_temp = in_addr_temp + 16*in_span; \
        m128_t11 = _mm_load_si128((__m128i *)in_addr_temp); \
         \
        m128_t0 = _mm_madd_epi16(m128_t8, m128_t0); \
        m128_t1 = _mm_madd_epi16(m128_t8, m128_t1); \
        m128_t1 = _mm_srli_si128(m128_t1,2); \
        m128_t1 = _mm_blend_epi16(m128_t0,m128_t1, 0x55); \
         \
         \
        m128_t2 = _mm_madd_epi16(m128_t9, m128_t2); \
        m128_t3 = _mm_madd_epi16(m128_t9, m128_t3); \
        m128_t3 = _mm_srli_si128(m128_t3,2); \
        m128_t3 = _mm_blend_epi16(m128_t2,m128_t3, 0x55); \
         \
         \
        m128_t4 = _mm_madd_epi16(m128_t10, m128_t4); \
        m128_t5 = _mm_madd_epi16(m128_t10, m128_t5); \
        m128_t5 = _mm_srli_si128(m128_t5,2); \
        m128_t5 = _mm_blend_epi16(m128_t4,m128_t5, 0x55); \
         \
         \
        m128_t6 = _mm_madd_epi16(m128_t11, m128_t6); \
        m128_t7 = _mm_madd_epi16(m128_t11, m128_t7); \
        m128_t7 = _mm_srli_si128(m128_t7,2); \
        m128_t7 = _mm_blend_epi16(m128_t6,m128_t7, 0x55); \
         \
        _MM_TRANSPOSE4_EPI32(m128_t1, m128_t3, m128_t5, m128_t7); \
         \
        out_addr_temp = OutBuf + ii*4; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t1);  out_addr_temp = out_addr_temp + 1; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t3);  out_addr_temp = out_addr_temp + 1; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t5);  out_addr_temp = out_addr_temp + 1; \
        _mm_store_si128((__m128i *)out_addr_temp, m128_t7); \
    } \
}

void fft2048(__m128i *InBuf, __m128i *OutBuf, __m128i *r32twiddle, __m128i *r64twiddle, __m128i *r2048twiddle)
{
    WORD32 i_loop=0;
    WORD32 in_span, out_span;

    __m128i m128_t0,m128_t1,m128_t2,m128_t3,m128_t4,m128_t5,m128_t6,m128_t7,m128_t8,m128_t9,m128_t10,m128_t11,m128_t12;
    __m128i OutTmp[512];
    __m128i fft2048_Temp64_32_Buf[64];

    __m128i  IQ_switch = _mm_setr_epi8(2,3,0,1,6,7,4,5,10,11,8,9,14,15,12,13);
	__m128i  Neg_I = _mm_setr_epi8(0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF);
	__m128i  Neg_R = _mm_setr_epi8(0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1, 0xFF, 0xFF, 0x1, 0x1);
	__m128i  Const_0707 = _mm_setr_epi8(0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A,0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A, 0x7F, 0x5A);
	__m128i  Const_0707_Minus = _mm_setr_epi8(0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5,0x81, 0xA5);


    in_span = 16;

    for (i_loop=0;i_loop<16;i_loop++)
    {
        radix32((__m128i *)(InBuf) + i_loop, in_span, (__m128i *)&fft2048_Temp64_32_Buf, (__m128i *)r32twiddle, (__m128i *)(OutTmp) + i_loop);
    }

    in_span = 1;
    out_span = 8;

    DIV_4_2((__m128i *)(OutTmp) + 0);
    DIV_4_2((__m128i *)(OutTmp) + 4);
    DIV_4_2((__m128i *)(OutTmp) + 8);
    DIV_4_2((__m128i *)(OutTmp) + 12);

    for (WORD32 i=1;i<4;i++)
    {
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 0, (__m128i *)r2048twiddle + (i*16 + 0) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 1, (__m128i *)r2048twiddle + (i*16 + 1) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 2, (__m128i *)r2048twiddle + (i*16 + 2) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 3, (__m128i *)r2048twiddle + (i*16 + 3) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 4, (__m128i *)r2048twiddle + (i*16 + 4) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 5, (__m128i *)r2048twiddle + (i*16 + 5) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 6, (__m128i *)r2048twiddle + (i*16 + 6) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 7, (__m128i *)r2048twiddle + (i*16 + 7) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 8, (__m128i *)r2048twiddle + (i*16 + 8) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 9, (__m128i *)r2048twiddle + (i*16 + 9) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 10,(__m128i *)r2048twiddle + (i*16 + 10) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 11, (__m128i *)r2048twiddle + (i*16 + 11) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 12, (__m128i *)r2048twiddle + (i*16 + 12) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 13, (__m128i *)r2048twiddle + (i*16 + 13) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 14,(__m128i *) r2048twiddle + (i*16 + 14) );
        MUL_CROSS((__m128i *)(OutTmp) + i*16 + 15, (__m128i *)r2048twiddle + (i*16 + 15) );
    }
    transpose64x4_mulzero((__m128i *)(OutTmp) + 0*64, in_span, (__m128i *)&fft2048_Temp64_32_Buf, r2048twiddle + 0*64);
    radix64((__m128i *)&fft2048_Temp64_32_Buf, in_span,(__m128i *)OutBuf + 0, out_span, (__m128i *)r64twiddle);

    for (i_loop=1;i_loop<8;i_loop++)
    {
        transpose64x4((__m128i *)(OutTmp) + i_loop*64, in_span, (__m128i *)&fft2048_Temp64_32_Buf, r2048twiddle + i_loop*64);
        radix64((__m128i *)&fft2048_Temp64_32_Buf, in_span,(__m128i *)OutBuf + i_loop, out_span, (__m128i *)r64twiddle);
    }
}

*/
import "C"
import "math/cmplx"
import "math"

type FFT2048Plan struct {
	r2048twiddle *Cmplx32v
	r32twiddle   *Cmplx32v
	r64twiddle   *Cmplx32v
}

func genTwiddle(stage0, stage1, s int) *Cmplx32v {
	var factor float64 = 32768.
	t := make([]complex128, stage0*stage1*s)
	for s0 := 0; s0 < stage0; s0 += 1 {
		for s1 := 0; s1 < stage1; s1 += 1 {
			for s2 := 0; s2 < s; s2 += 1 {
				p := -1. * math.Pi * 2. * float64(s0) * float64(s1) / float64(stage0*stage1)
				a := complex(factor, 0.) * cmplx.Exp(complex(0., p))
				t[s0*stage1*s+s1*s+s2] = complex(imag(a), real(a))
			}
		}
	}
	return ToM128Buf(t)
}

func NewFFT2048Plan() *FFT2048Plan {
	r := &FFT2048Plan{}
	r.r2048twiddle = genTwiddle(32, 64, 1)
	r.r32twiddle = genTwiddle(8, 4, 4)
	r.r64twiddle = genTwiddle(8, 8, 4)
	return r
}

func (self *FFT2048Plan) Do(in, out *Cmplx32v) {
	C.fft2048((*C.__m128i)(in.d), (*C.__m128i)(out.d), (*C.__m128i)(self.r32twiddle.d), (*C.__m128i)(self.r64twiddle.d), (*C.__m128i)(self.r2048twiddle.d))
}
