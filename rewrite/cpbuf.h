#include "spin_lock.h"
#define METASIZE 256
#define CAPSIZE 256
typedef long long int64;
typedef unsigned long long uint64;
typedef int int32;
typedef unsigned int uint32;

struct structCPBLen {
  int64 resLen;
  int64 dataLen;
  int64 cpLen;
  int64 version;
  int64 type;
};
typedef struct CPBMeta {
  char name[METASIZE-sizeof(struct structCPBLen)-sizeof(char)*16];
  struct structCPBLen cpbLen;
  char key[16];
} CPBMeta;
typedef struct {
	CPBMeta meta;
  raw_spinlock_t lockM;
  int64 _brk;
} VMemHead;
typedef struct {
  CPBMeta meta;
  raw_spinlock_t lock;
  char buf[256];
} TimingBufHead;

typedef struct {
  CPBMeta meta;
  char buf[256];
} GBufHead;
typedef struct {
  int64 _brk;
  int64 _free;
  int64 _wr;
  int64 _rd;
} lmmem;
typedef struct {
  int32 _brk;
  int32 _free;
  int32 _wr;
  int32 _rd;
} smmem;
typedef struct {
  int64 cpuoff;
  int64 sysoff;
  float64 cpua;
  float64 cpuo;
} clockTrans;

typedef struct {
  char name[32];
  char entity[32];
  int32 version;
  int32 seqnum;
  int32 resoff;//from Res[0]
  int32 reslen;//for check
  int32 dummy; 
  char typemd5[16];
  char entitymd5[16];
  raw_spinlock_t lock[4];
  lmmem volatileMem;
  clockTrans clkt;
  smmem preAllocMem;
  char res[64];

} capability;