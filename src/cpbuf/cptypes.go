package cpbuf
/*
typedef char int8;
typedef unsigned char uint8;
typedef short int16;
typedef unsigned short uint16;
typedef int int32;
typedef long long int int64;
typedef unsigned int uint32;
typedef unsigned long long int uint64;
typedef float float32;
typedef double float64;
typedef unsigned char byte;
// Version 4
typedef struct Version_s {
	uint16 build;
	uint8 minor;
	uint8 major;
} Version;
// ResBlk 8
typedef struct ResBlk_s {
	int32 offset;
	int32 len;
} ResBlk;
// SMMem 16
typedef struct SMMem_s {
	int32 _brk;
	int32 _free;
	int32 _wr;
	int32 _rd;
} SMMem;
// raw_spin_lock_t 4
typedef struct raw_spin_lock_t_s {
	uint32 lock;
} raw_spin_lock_t;
// ClkTrans 32
typedef struct ClkTrans_s {
	int64 cpuoff;
	int64 sysoff;
	float64 clkr;
	float64 clks;
} ClkTrans;
// NumInApp 4
typedef struct NumInApp_s {
	uint16 num;
	uint8 fun;
	uint8 app;
} NumInApp;
// CpInfo 24
typedef struct CpInfo_s {
	int64 resLen;
	int64 dataLen;
	int64 cpLen;
} CpInfo;
// LMMem 32
typedef struct LMMem_s {
	int64 _brk;
	int64 _free;
	int64 _wr;
	int64 _rd;
} LMMem;
// CapMeta 72
typedef struct CapMeta_s {
	char name[32];
	char entity[32];
	Version ver;
	NumInApp num;
} CapMeta;
// CapDefault 96
typedef struct CapDefault_s {
	raw_spin_lock_t lock[4];
	LMMem volatileMem;
	ClkTrans clk;
	SMMem preAllocMem;
} CapDefault;
// CpMeta 96
typedef struct CpMeta_s {
	char name[32];
	char app[32];
	CpInfo info;
	Version ver;
	int32 salt;
} CpMeta;
// IndexItem 32
typedef struct IndexItem_s {
	ResBlk blk;
	int64 captype;
	byte md5[16];
} IndexItem;
// IndexHead 1024
typedef struct IndexHead_s {
	IndexItem index[32];
} IndexHead;
// Capability 2048
typedef struct Capability_s {
	CapMeta meta;
	CapDefault cap;
	byte userdef[1880];
} Capability;
// CPBuffer 1024
typedef struct CPBuffer_s {
	CpMeta meta;
	byte userdef[928];
} CPBuffer;
// BufHead 65536
typedef struct BufHead_s {
	IndexHead index;
	CPBuffer cpbuf;
	Capability Caps[31];
} BufHead;
// ResMem 1048576
typedef struct ResMem_s {
	BufHead head;
	byte userdef[983040];
} ResMem;
*/
import "C"
// Version 4
type Version struct {
	build uint16
	minor uint8
	major uint8
}
// ResBlk 8
type ResBlk struct {
	offset int32
	len int32
}
// SMMem 16
type SMMem struct {
	_brk int32
	_free int32
	_wr int32
	_rd int32
}
// raw_spin_lock_t 4
type raw_spin_lock_t struct {
	lock uint32
}
// ClkTrans 32
type ClkTrans struct {
	cpuoff int64
	sysoff int64
	clkr float64
	clks float64
}
// NumInApp 4
type NumInApp struct {
	num uint16
	fun uint8
	app uint8
}
// CpInfo 24
type CpInfo struct {
	resLen int64
	dataLen int64
	cpLen int64
}
// LMMem 32
type LMMem struct {
	_brk int64
	_free int64
	_wr int64
	_rd int64
}
// CapMeta 72
type CapMeta struct {
	name [32]byte
	entity [32]byte
	ver Version
	num NumInApp
}
// CapDefault 96
type CapDefault struct {
	lock [4]raw_spin_lock_t
	volatileMem LMMem
	clk ClkTrans
	preAllocMem SMMem
}
// CpMeta 96
type CpMeta struct {
	name [32]byte
	app [32]byte
	info CpInfo
	ver Version
	salt int32
}
// IndexItem 32
type IndexItem struct {
	blk ResBlk
	captype int64
	md5 [16]byte
}
// IndexHead 1024
type IndexHead struct {
	index [32]IndexItem
}
// Capability 2048
type Capability struct {
	meta CapMeta
	cap CapDefault
	userdef [1880]byte
}
// CPBuffer 1024
type CPBuffer struct {
	meta CpMeta
	userdef [928]byte
}
// BufHead 65536
type BufHead struct {
	index IndexHead
	cpbuf CPBuffer
	Caps [31]Capability
}
// ResMem 1048576
type ResMem struct {
	head BufHead
	userdef [983040]byte
}
