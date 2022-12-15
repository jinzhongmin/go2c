package dlfcn

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
void *rtdl_default(){ return RTLD_DEFAULT; }
void *rtdl_next(){ return RTLD_NEXT; }
*/
import "C"
import "unsafe"

var RTLD_DEFAULT unsafe.Pointer
var RTLD_NEXT unsafe.Pointer

func init() {
	RTLD_DEFAULT = C.rtdl_default()
	RTLD_NEXT = C.rtdl_next()
}

type Mode C.int

const (
	RTLDNow    Mode = C.RTLD_NOW
	RTLDLazy   Mode = C.RTLD_LAZY
	RTLDGlobal Mode = C.RTLD_GLOBAL
	RTLDLocal  Mode = C.RTLD_LOCAL
)

func (m Mode) toC() C.int {
	return C.int(m)
}
