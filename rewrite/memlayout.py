
class type:
	def __init__(self,name,size,ctypename,gotypename='default'):
		self.name = name
		self.size = size
		self.ctypename = ctypename
		if gotypename=='default':
			self.gotypename = name
		else:
			self.gotypename = gotypename

	def cgen(self,v):
		print self.name,'%s;'%v.name
	def gogen(self,v):
		print v.name,self.gotypename
	def carray(self,v):
		print self.name,'%s[%d];'%(v.name,v.l)
	def goarray(self,v):
		print v.name,'[%d]%s'%(v.l,self.gotypename)
		
class struct:
	def __init__(self,types,name):
		self.types = types
		self.name = name
		self.items = []
		self.size = 0
	def fromList(self,list):
		l = 0
		self.m = genMap(self.types)
		for x in list:
			if x.name!='userdef':
				self.items.append(x)
				if self.m[x.type].size>=4 and (l%4)!=0:
					print "//",x.name,"align to",l	
				l += self.m[x.type].size*x.l
			else:
				self.items.append(item(x.name,'byte',x.l-l))
				l = x.l
		self.types.append(type(self.name,l,"struct %s_s"%self.name))
		self.size = l
		return self.types
	def ctypeDec(self):
		print '//',self.name,self.size
		print 'typedef','struct',self.name+'_s','{'
		for x in self.items:
			print '\t',
			if x.l==1:
				self.m[x.type].cgen(x)
			else:
				self.m[x.type].carray(x)
		print '}','%s;'%self.name

	def gotypeDec(self):
		print '//',self.name,self.size
		print 'type',self.name,'struct','{'
		for x in self.items:
			print '\t',
			if x.l==1:
				self.m[x.type].gogen(x)
			else:
				self.m[x.type].goarray(x)
		print '}'


def genTypes():
	types = []
	types.append(type('int8',1,'char'))
	types.append(type('uint8',1,'unsigned char'))
	types.append(type('int16',2,'short'))
	types.append(type('uint16',2,'unsigned short'))
	types.append(type('int32',4,'int'))
	types.append(type('int64',8,'long long int'))
	types.append(type('uint32',4,'unsigned int'))
	types.append(type('uint64',8,'unsigned long long int'))
	types.append(type('float32',4,'float','uint32'))
	types.append(type('float64',8,'double'))
	types.append(type('byte',1,'unsigned char'))
	types.append(type('char',1,'char','byte'))
	return types

def ctypeDec(types):
	for x in types:
		if x.name!=x.ctypename:
			print 'typedef',x.ctypename,'%s;'%x.name

def genMap(types):
	m = {}
	for x in types:
		m[x.name]=x
	return m
class item:
	def __init__(self,name,t,l):
		self.name = name
		self.type = t
		self.l =l
def str2list(str):
	l = []
	for line in str.split('\n'):
		line = line.replace('\n','')
		it = line.split(' ')
		if len(it)!=3:
			continue
		l.append(item(it[0],it[1],int(it[2])))
	return l
raw_spin_lock_t_str = """
lock uint32 1
"""
CpInfo_str = """
resLen int64 1
dataLen int64 1
cpLen int64 1
"""
Version_str = """
build uint16 1
minor uint8 1
major uint8 1
"""
NumInApp_str = """
num uint16 1
fun uint8 1
app uint8 1
"""
CpMeta_str = """
name char 32
app char 32
info CpInfo 1
ver Version 1
salt int32 1
"""
ResBlk_str = """
offset int32 1
len int32 1
"""
IndexItem_str = """
blk ResBlk 1
captype int64 1
md5 byte 16
"""
IndexHead_str = """
index IndexItem 32
"""
CPBuffer_str = """
meta CpMeta 1
userdef byte 1024
"""
SMMem_str = """
_brk int32 1
_free int32 1
_wr int32 1
_rd int32 1
"""
LMMem_str = """
_brk int64 1
_free int64 1
_wr int64 1
_rd int64 1
"""
ClkTrans_str = """
cpuoff int64 1
sysoff int64 1
clkr float64 1
clks float64 1
"""
CapMeta_str = """
name char 32
entity char 32
blk ResBlk
ver Version 1
num NumInApp 1
"""
CapDefault_str = """
lock raw_spin_lock_t 4
volatileMem LMMem 1
clk ClkTrans 1
preAllocMem SMMem 1
"""
Capability_str = """
meta CapMeta 1
cap CapDefault 1
userdef byte 2048
"""
BufHead_str = """
index IndexHead 1
cpbuf CPBuffer 1
Caps Capability 31
"""
ResMem_str = """
head BufHead 1
userdef byte 1048576
"""
structs = [   'raw_spin_lock_t'
			, 'CpInfo'
			, 'Version'
			, 'CpMeta'
			, 'IndexItem'
			, 'Sector00' 
			, 'Sector01' 
			]
def genStructList(lmap):
	s = []
	for k in lmap:
		if k[-4:]=='_str':
			s.append(k[:-4])
	m = {}
	for x in s:
		m[x]=[]
		for y in s:
			if y!=x:
				if y in lmap[x+'_str']:
					m[x].append(y)
	ss = []
	while m != {}:
		for k in m:
			if len(m[k])==0:
				ss.append(k)
				
		for x in ss:
			if x in m:
				del m[x]
			for k in m:
				if x in m[k]:
					m[k].remove(x)
		#print ss
		#print m

	return ss
def cDec(structs,lmap):
	types = genTypes()
	for n in structs:
		struct_item = struct(types,n)
		s = lmap[n+'_str']
		types = struct_item.fromList(str2list(s))
		struct_item.ctypeDec()

def goDec(structs,lmap):
	types = genTypes()
	for n in structs:
		struct_item = struct(types,n)
		s = lmap[n+'_str']
		types = struct_item.fromList(str2list(s))
		struct_item.gotypeDec()

if __name__=='__main__':
	print "package cpbuf"
	print "/*"
	types = genTypes()
	ctypeDec(types)
	s = genStructList(locals())
	cDec(s,locals())
	print "*/"
	print 'import "C"'
	goDec(s,locals())
	