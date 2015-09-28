#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <string.h>

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
typedef char WORD8;
#define MUL_CROSS(in_addr, twiddle_addr) \
{\
	__m128i * twiddleaddr_temp;\
	__m128i * in_addr_temp;\
	in_addr_temp = in_addr;\
	m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp));\
	twiddleaddr_temp = twiddle_addr;\
	m128_t3 = _mm_load_si128((__m128i *)twiddleaddr_temp); \
	twiddleaddr_temp = twiddleaddr_temp + 1;\
	m128_t4 = _mm_load_si128((__m128i *)twiddleaddr_temp); \
	m128_t1 = _mm_madd_epi16(m128_t0, m128_t3);  \
	m128_t0 = _mm_madd_epi16(m128_t0, m128_t4); \
	m128_t0 = _mm_srli_si128(m128_t0,2);\
	m128_t1 = _mm_blend_epi16(m128_t1,m128_t0, 0x55); \
	_mm_store_si128((__m128i *)in_addr_temp, m128_t1);  \
}

#define DIV_4_2(in_addr) \
{\
	__m128i * in_addr_temp;\
	in_addr_temp = in_addr;\
	m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
	_mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0,1));  \
	in_addr_temp = in_addr_temp + 1;\
	m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
	_mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0,1));  \
	in_addr_temp = in_addr_temp + 1;\
	m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
	_mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0,1));  \
	in_addr_temp = in_addr_temp + 1;\
	m128_t0 = _mm_load_si128((__m128i *)(in_addr_temp)); \
	_mm_store_si128((__m128i *)in_addr_temp, _mm_srai_epi16(m128_t0,1));  \
	in_addr_temp = in_addr_temp + 1;\
}

#define radix4_register(in0,in1,in2,in3,out_addr,out_span) \
{\
	__m128i * in_addr_temp;\
	m128_t2 = _mm_adds_epi16(in0,in2); /*AA*/\
	m128_t6 = _mm_adds_epi16(in1,in3); /*CC*/\
	m128_t7=  _mm_adds_epi16(m128_t2,m128_t6);  /*out 1*/ \
	m128_t2 = _mm_subs_epi16(m128_t2,m128_t6);  /*out 3*/ \
	m128_t6 =  m128_t7; /*out 1*/\
	\
	m128_t3 = _mm_subs_epi16(in0,in2); /*BB*/\
	m128_t12 = _mm_subs_epi16(in1,in3); /*DD*/\
	m128_t12 = _mm_shuffle_epi8(m128_t12, IQ_switch);\
	m128_t12 = _mm_sign_epi16(m128_t12, Neg_R);/*j*D*/\
	\
	m128_t7 = _mm_adds_epi16(m128_t3,m128_t12); /*out 2 */\
	m128_t3 = _mm_subs_epi16(m128_t3,m128_t12); /*out 4 */\
	/*m128_t6  m128_t7 m128_t2  m128_t3*/ \
	\
	_mm_store_si128((__m128i *)out_addr, m128_t6);  \
	out_addr = out_addr + out_span;\
	\
	_mm_store_si128((__m128i *)out_addr, m128_t7);  \
	out_addr = out_addr + out_span;\
	\
	_mm_store_si128((__m128i *)out_addr, m128_t2);  \
	out_addr = out_addr + out_span;\
	\
	_mm_store_si128((__m128i *)out_addr, m128_t3);  \
	\
}

#define radix8_0(in_addr,in_span,out_addr,out_span)\
{\
	__m128i * out_addr_temp1,*out_addr_temp2;\
	__m128i * in_addr_temp;\
	in_addr_temp = in_addr;\
	m128_t0 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span;\
	m128_t1 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span;\
	m128_t2 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span;\
	m128_t3 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span;\
	m128_t4 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span;\
	m128_t5 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span;\
	m128_t6 = _mm_load_si128((__m128i *)in_addr_temp);  in_addr_temp = in_addr_temp + in_span;\
	m128_t7 = _mm_load_si128((__m128i *)in_addr_temp);\
	\
	m128_t8 = _mm_adds_epi16(m128_t0, m128_t4); /* Temp8_0 */ \
	m128_t9 = _mm_subs_epi16(m128_t0, m128_t4); /* Temp8_1 */ \
	m128_t10 = _mm_adds_epi16(m128_t1, m128_t5);/* Temp8_2 */ \
	m128_t11 = _mm_subs_epi16(m128_t1, m128_t5);/* Temp8_3 */ \
	m128_t0 = _mm_adds_epi16(m128_t2, m128_t6); /* Temp8_4 */ \
	m128_t1 = _mm_subs_epi16(m128_t2, m128_t6);/* Temp8_5 */ \
	m128_t4 = _mm_adds_epi16(m128_t3, m128_t7);/* Temp8_6 */ \
	m128_t5 = _mm_subs_epi16(m128_t3, m128_t7);/* Temp8_7 */ \
	\
	/*x1(1+5)  = x1(1 +5) * (0 - i) = imag - real*j */\
	m128_t1 = _mm_shuffle_epi8(m128_t1, IQ_switch);\
	m128_t1 = _mm_sign_epi16(m128_t1, Neg_R); /* Temp8_5 */ \
	\
	m128_t2 = _mm_mulhrs_epi16(m128_t11, Const_0707); /*Temp8_3*/ \
	m128_t3 = _mm_mulhrs_epi16(m128_t5, Const_0707_Minus);  /*Temp8_7*/ \
	\
	m128_t6 = _mm_shuffle_epi8(m128_t2, IQ_switch); /*Temp8_3*/ \
	m128_t7 = _mm_adds_epi16(m128_t2,m128_t6); \
	m128_t2 = _mm_subs_epi16(m128_t2,m128_t6);\
	m128_t11 = _mm_blend_epi16(m128_t2,m128_t7, 0xAA);\
	\
	m128_t6 = _mm_shuffle_epi8(m128_t3, IQ_switch); /*Temp8_7*/ \
	m128_t7 = _mm_adds_epi16(m128_t3,m128_t6); \
	m128_t3 = _mm_subs_epi16(m128_t3,m128_t6);\
	m128_t5 = _mm_blend_epi16(m128_t3,m128_t7, 0x55);\
	\
	\
	out_addr_temp1 = out_addr;\
	out_addr_temp2 = out_addr + out_span;\
	radix4_register(m128_t8,m128_t10,m128_t0,m128_t4, out_addr_temp1, 2*out_span);\
	radix4_register(m128_t9,m128_t11,m128_t1,m128_t5, out_addr_temp2, 2*out_span);\
	\
}

#define radix4_register_mul_zero(in0, in1, in2, in3, out_addr, out_span, twiddle_addr) \
{ \
	__m128i *in_addr_temp; \
	__m128i *twiddle_addr_temp; \
	m128_t2 = _mm_adds_epi16(in0,in2); /*AA*/ \
	m128_t6 = _mm_adds_epi16(in1,in3); /*CC*/ \
	m128_t7=  _mm_adds_epi16(m128_t2,m128_t6);  /*out 1*/ \
	m128_t2 = _mm_subs_epi16(m128_t2,m128_t6);  /*out 3*/ \
	m128_t6 =  m128_t7; /*out 1*/ \
	 \
	m128_t3 = _mm_subs_epi16(in0,in2); /*BB*/ \
	m128_t12 = _mm_subs_epi16(in1,in3); /*DD*/ \
	m128_t12 = _mm_shuffle_epi8(m128_t12, IQ_switch); \
	m128_t12 = _mm_sign_epi16(m128_t12, Neg_R);/*j*D*/ \
	 \
	m128_t7 = _mm_adds_epi16(m128_t3,m128_t12); /*out 2 */ \
	m128_t3 = _mm_subs_epi16(m128_t3,m128_t12); /*out 4 */ \
	/*m128_t6  m128_t7 m128_t2  m128_t3*/ \
	 \
	 \
	_mm_store_si128((__m128i *)out_addr,  _mm_srai_epi16(m128_t6,1));   out_addr = out_addr + out_span; \
	 \
	twiddle_addr_temp = twiddle_addr + 2 * 2; \
	m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 3; \
	m128_t0 = _mm_madd_epi16(m128_t7, m128_t0); \
	m128_t4 = _mm_madd_epi16(m128_t7, m128_t4); \
	m128_t4 = _mm_srli_si128(m128_t4,2); \
	m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
	_mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
	 \
	m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 3; \
	m128_t0 = _mm_madd_epi16(m128_t2, m128_t0); \
	m128_t4 = _mm_madd_epi16(m128_t2, m128_t4); \
	m128_t4 = _mm_srli_si128(m128_t4,2); \
	m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
	_mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
	 \
	m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
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
	m128_t2 = _mm_adds_epi16(in0,in2); /*AA*/ \
	m128_t6 = _mm_adds_epi16(in1,in3); /*CC*/ \
	m128_t7=  _mm_adds_epi16(m128_t2,m128_t6);  /*out 1*/ \
	m128_t2 = _mm_subs_epi16(m128_t2,m128_t6);  /*out 3*/ \
	m128_t6 =  m128_t7; /*out 1*/ \
	 \
	m128_t3 = _mm_subs_epi16(in0,in2); /*BB*/ \
	m128_t12 = _mm_subs_epi16(in1,in3); /*DD*/ \
	m128_t12 = _mm_shuffle_epi8(m128_t12, IQ_switch); \
	m128_t12 = _mm_sign_epi16(m128_t12, Neg_R);/*j*D*/ \
	 \
	m128_t7 = _mm_adds_epi16(m128_t3,m128_t12); /*out 2 */ \
	m128_t3 = _mm_subs_epi16(m128_t3,m128_t12); /*out 4 */ \
	/*m128_t6  m128_t7 m128_t2  m128_t3*/ \
	 \
	twiddle_addr_temp = twiddle_addr; \
	m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 3; \
	m128_t0 = _mm_madd_epi16(m128_t6, m128_t0); \
	m128_t4 = _mm_madd_epi16(m128_t6, m128_t4); \
	m128_t4 = _mm_srli_si128(m128_t4,2); \
	m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
	_mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
	 \
	m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 3; \
	m128_t0 = _mm_madd_epi16(m128_t7, m128_t0); \
	m128_t4 = _mm_madd_epi16(m128_t7, m128_t4); \
	m128_t4 = _mm_srli_si128(m128_t4,2); \
	m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
	_mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
	 \
	m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 3; \
	m128_t0 = _mm_madd_epi16(m128_t2, m128_t0); \
	m128_t4 = _mm_madd_epi16(m128_t2, m128_t4); \
	m128_t4 = _mm_srli_si128(m128_t4,2); \
	m128_t0 = _mm_blend_epi16(m128_t0,m128_t4, 0x55); \
	_mm_store_si128((__m128i *)out_addr, m128_t0);   out_addr = out_addr + out_span; \
	 \
	m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
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
	m128_t8 = _mm_adds_epi16(m128_t0, m128_t4); /* Temp8_0 */ \
	m128_t9 = _mm_subs_epi16(m128_t0, m128_t4); /* Temp8_1 */ \
	m128_t10 = _mm_adds_epi16(m128_t1, m128_t5);/* Temp8_2 */ \
	m128_t11 = _mm_subs_epi16(m128_t1, m128_t5);/* Temp8_3 */ \
	m128_t0 = _mm_adds_epi16(m128_t2, m128_t6); /* Temp8_4 */ \
	m128_t1 = _mm_subs_epi16(m128_t2, m128_t6);/* Temp8_5 */ \
	m128_t4 = _mm_adds_epi16(m128_t3, m128_t7);/* Temp8_6 */ \
	m128_t5 = _mm_subs_epi16(m128_t3, m128_t7);/* Temp8_7 */ \
	 \
	/*x1(1+5)  = x1(1 +5) * (0 - i) = imag - real*j */ \
	m128_t1 = _mm_shuffle_epi8(m128_t1, IQ_switch); \
	m128_t1 = _mm_sign_epi16(m128_t1, Neg_R); /* Temp8_5 */ \
	 \
	m128_t2 = _mm_mulhrs_epi16(m128_t11, Const_0707); /*Temp8_3*/ \
	m128_t3 = _mm_mulhrs_epi16(m128_t5, Const_0707_Minus);  /*Temp8_7*/ \
	 \
	m128_t6 = _mm_shuffle_epi8(m128_t2, IQ_switch); /*Temp8_3*/ \
	m128_t7 = _mm_adds_epi16(m128_t2,m128_t6); \
	m128_t2 = _mm_subs_epi16(m128_t2,m128_t6); \
	m128_t11 = _mm_blend_epi16(m128_t2,m128_t7, 0xAA); \
	 \
	m128_t6 = _mm_shuffle_epi8(m128_t3, IQ_switch); /*Temp8_7*/ \
	m128_t7 = _mm_adds_epi16(m128_t3,m128_t6); \
	m128_t3 = _mm_subs_epi16(m128_t3,m128_t6); \
	m128_t5 = _mm_blend_epi16(m128_t3,m128_t7, 0x55); \
	 \
	 \
	out_addr_temp1 = out_addr; \
	out_addr_temp2 = out_addr + out_span; \
	radix4_register_mul_zero(m128_t8,m128_t10,m128_t0,m128_t4, out_addr_temp1, 2*out_span,twiddle_addr); \
	radix4_register_mul(m128_t9,m128_t11,m128_t1,m128_t5, out_addr_temp2, 2*out_span,twiddle_addr + 2); \
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
	m128_t4 = _mm_adds_epi16(m128_t0,m128_t2); /*AA*/ \
	m128_t5 = _mm_adds_epi16(m128_t1,m128_t3); /*CC*/ \
	m128_t6 = _mm_subs_epi16(m128_t0,m128_t2); /*BB*/ \
	m128_t7 = _mm_subs_epi16(m128_t1,m128_t3); /*DD*/ \
	m128_t7 = _mm_shuffle_epi8(m128_t7, IQ_switch); \
	m128_t7 = _mm_sign_epi16(m128_t7, Neg_R);/*j*D*/ \
	 \
	m128_t8 = _mm_adds_epi16(m128_t4,m128_t5); \
	m128_t8 = _mm_srai_epi16(m128_t8,1); \
	m128_t9 = _mm_adds_epi16(m128_t6,m128_t7); \
	m128_t9 = _mm_srai_epi16(m128_t9,1); \
	m128_t10 = _mm_subs_epi16(m128_t4,m128_t5); \
	m128_t10 = _mm_srai_epi16(m128_t10,1); \
	m128_t11 = _mm_subs_epi16(m128_t6,m128_t7); \
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
	m128_t4 = _mm_adds_epi16(m128_t0,m128_t2); /*AA*/ \
	m128_t5 = _mm_adds_epi16(m128_t1,m128_t3); /*CC*/ \
	m128_t6 = _mm_subs_epi16(m128_t0,m128_t2); /*BB*/ \
	m128_t7 = _mm_subs_epi16(m128_t1,m128_t3); /*DD*/ \
	m128_t7 = _mm_shuffle_epi8(m128_t7, IQ_switch); \
	m128_t7 = _mm_sign_epi16(m128_t7, Neg_R);/*j*D*/ \
	 \
	m128_t8 = _mm_adds_epi16(m128_t4,m128_t5); \
	m128_t9 = _mm_adds_epi16(m128_t6,m128_t7); \
	m128_t10 = _mm_subs_epi16(m128_t4,m128_t5); \
	m128_t11 = _mm_subs_epi16(m128_t6,m128_t7); \
	twiddle_addr_temp = twiddle_addr + 2; \
	m128_t2 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t3 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t5 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t6 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
	m128_t7 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
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

#define transpose64x4_mulzero(InBuf, in_span, OutBuf)\
{\
	__m128i *in_addr_temp; \
	__m128i *out_addr_temp; \
    WORD32 in_span_temp = 16 * in_span; \
	 \
	for (WORD32 ii=0;ii<64/4;ii++) \
	{ \
		in_addr_temp = InBuf + ii*in_span; \
		m128_t1 = _mm_load_si128(in_addr_temp); in_addr_temp = in_addr_temp + in_span_temp; \
		m128_t3 = _mm_load_si128(in_addr_temp); in_addr_temp = in_addr_temp + in_span_temp; \
		m128_t5 = _mm_load_si128(in_addr_temp); in_addr_temp = in_addr_temp + in_span_temp; \
		m128_t7 = _mm_load_si128(in_addr_temp); \
		 \
		_MM_TRANSPOSE4_EPI32(m128_t1, m128_t3, m128_t5, m128_t7); \
		 \
		out_addr_temp = OutBuf + ii*4; \
		_mm_store_si128(out_addr_temp, m128_t1);  out_addr_temp = out_addr_temp + 1; \
		_mm_store_si128(out_addr_temp, m128_t3);  out_addr_temp = out_addr_temp + 1; \
		_mm_store_si128(out_addr_temp, m128_t5);  out_addr_temp = out_addr_temp + 1; \
		_mm_store_si128(out_addr_temp, m128_t7); \
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
		twiddle_addr_temp = twiddle_addr + ii * 2; \
		m128_t0 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
		m128_t1 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 31; \
		m128_t2 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
		m128_t3 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 31; \
		m128_t4 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
		m128_t5 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 31; \
		m128_t6 = _mm_load_si128((__m128i *)twiddle_addr_temp); twiddle_addr_temp = twiddle_addr_temp + 1; \
		m128_t7 = _mm_load_si128((__m128i *)twiddle_addr_temp); \
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

#define radix64(InBuf ,in_span, OutBuf , out_span, ifft2048_r64twiddle)\
{\
	radix8_0(InBuf + 0*in_span, 8*in_span, ifft2048_Temp64_32_Buf +0*8,1);\
	radix8_0_mul(InBuf + 1*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 1*8, 1, ifft2048_r64twiddle + 2*4*1*2); \
	radix8_0_mul(InBuf + 2*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 2*8, 1, ifft2048_r64twiddle + 2*4*2*2); \
	radix8_0_mul(InBuf + 3*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 3*8, 1, ifft2048_r64twiddle + 2*4*3*2); \
	radix8_0_mul(InBuf + 4*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 4*8, 1, ifft2048_r64twiddle + 2*4*4*2); \
	radix8_0_mul(InBuf + 5*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 5*8, 1, ifft2048_r64twiddle + 2*4*5*2); \
	radix8_0_mul(InBuf + 6*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 6*8, 1, ifft2048_r64twiddle + 2*4*6*2); \
	radix8_0_mul(InBuf + 7*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 7*8, 1, ifft2048_r64twiddle + 2*4*7*2); \
	\
	DIV_4_2((__m128i *)(ifft2048_Temp64_32_Buf) + 0*4);\
	DIV_4_2((__m128i *)(ifft2048_Temp64_32_Buf) + 1*4);\
	\
	radix8_0(ifft2048_Temp64_32_Buf + 0, 8, OutBuf + 0*out_span, 8*out_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 1, 8, OutBuf + 1*out_span, 8*out_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 2, 8, OutBuf + 2*out_span, 8*out_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 3, 8, OutBuf + 3*out_span, 8*out_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 4, 8, OutBuf + 4*out_span, 8*out_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 5, 8, OutBuf + 5*out_span, 8*out_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 6, 8, OutBuf + 6*out_span, 8*out_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 7, 8, OutBuf + 7*out_span, 8*out_span);\
	\
}


#define radix32(InBuf,  in_span,  ifft2048_Temp64_32_Buf,  ifft2048_r32twiddle, OutBuf)\
{\
	radix4_0_zeromul(InBuf + 0*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 0*4); \
	radix4_0_mul(InBuf + 1*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 1*4, ifft2048_r32twiddle + 2*4*1); \
	radix4_0_mul(InBuf + 2*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 2*4, ifft2048_r32twiddle + 2*4*2); \
	radix4_0_mul(InBuf + 3*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 3*4, ifft2048_r32twiddle + 2*4*3); \
	radix4_0_mul(InBuf + 4*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 4*4, ifft2048_r32twiddle + 2*4*4); \
	radix4_0_mul(InBuf + 5*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 5*4, ifft2048_r32twiddle + 2*4*5); \
	radix4_0_mul(InBuf + 6*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 6*4, ifft2048_r32twiddle + 2*4*6); \
	radix4_0_mul(InBuf + 7*in_span, 8*in_span, ifft2048_Temp64_32_Buf + 7*4, ifft2048_r32twiddle + 2*4*7); \
	\
	radix8_0(ifft2048_Temp64_32_Buf + 0, 4, OutBuf + 0*in_span, 4*in_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 1, 4, OutBuf + 1*in_span, 4*in_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 2, 4, OutBuf + 2*in_span, 4*in_span);\
	radix8_0(ifft2048_Temp64_32_Buf + 3, 4, OutBuf + 3*in_span, 4*in_span);\
}



