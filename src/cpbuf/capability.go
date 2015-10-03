package cpbuf

type BufCapability interface {
	name() string
	attach(*GBuf)
	key() string
	concurrent() interface{}
	lock()
	unlcok()
}

type SharedObj struct {
	Obj    interface{}
	lock   func()
	unlock func()
}

func key() string {

}
