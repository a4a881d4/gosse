package sse

/*
#include<stdio.h>
#include<stdlib.h>
#include"align.h"


void *align_malloc( void** praw, int l, int a){
	void *raw = malloc(l+a-1);
	void *r = (void *)(((unsigned long)raw+a-1)&(unsigned long)(-a));
	*praw = raw;
	printf("raw=%p,align=%p\n",raw,r);
	return r;
}
*/
import "C"
