package dlfcn

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/jinzhongmin/go2c/pkg/mem"
)

var iconv func(string) string

func SetErrIconv(fn func(string) string) {
	iconv = fn
}

func DlError() error {
	e := C.dlerror()
	if e == nil {
		return nil
	}
	b := C.GoString(e)
	if iconv != nil {
		s := iconv(b)
		return errors.New(s)
	}
	return errors.New(b)
}

type Handle struct {
	c unsafe.Pointer
}

func DlOpen(file string, mod Mode) (*Handle, error) {
	f := C.CString(file)
	defer mem.Free(unsafe.Pointer(f))

	h := C.dlopen(f, mod.toC())
	if h == nil {
		return nil, DlError()
	}
	di := new(Handle)
	di.c = h
	return di, nil
}
func (hd *Handle) Close() {
	if hd.c != nil {
		C.dlclose(hd.c)
	}
}
func (hd Handle) Symbol(name string) (unsafe.Pointer, error) {
	n := C.CString(name)
	defer mem.Free(unsafe.Pointer(n))

	p := C.dlsym(hd.c, n)
	if p == nil {
		return nil, DlError()
	}

	return p, nil
}

//warp dladdr
func DlAddr(addr unsafe.Pointer) (int, string, unsafe.Pointer, string, unsafe.Pointer) {
	di := new(C.Dl_info)
	i := C.dladdr(addr, di)
	fname := C.GoString(di.dli_fname)
	fbase := di.dli_fbase
	sname := C.GoString(di.dli_sname)
	saddr := di.dli_saddr
	return int(int32(i)), fname, fbase, sname, saddr
}

//warp dlsym
func DlSym(p unsafe.Pointer, name string) (unsafe.Pointer, error) {
	n := C.CString(name)
	defer mem.Free(unsafe.Pointer(n))
	r := C.dlsym(p, n)
	if r == nil {
		return nil, DlError()
	}
	return r, nil
}
