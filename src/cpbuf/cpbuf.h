#include "spin_lock.h"
#define METASIZE 256
typedef long long int64;
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

