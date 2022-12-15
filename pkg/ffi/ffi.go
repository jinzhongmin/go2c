package ffi

/*
#include <ffi.h>
extern void closure_caller(ffi_cif* cif, void* ret, void* args, void* user_data);
*/
import "C"
import (
	"reflect"
	"unsafe"

	"github.com/jinzhongmin/go2c/pkg/mem"
)

type Cif struct {
	cif    *C.ffi_cif
	retVal unsafe.Pointer
}

func NewCif(abi Abi, ret *Type, args ...*Type) (*Cif, error) {
	cif := new(Cif)

	cifLen := mem.Sizeof(C.ffi_cif{})
	cif.cif = (*C.ffi_cif)(mem.Malloc(1, cifLen))
	mem.Memset(unsafe.Pointer(cif.cif), 0, cifLen)

	//malloc return val ptr
	if ret == nil {
		cif.retVal = mem.Malloc(1, int(TypeVoid.size))
	} else {
		cif.retVal = mem.Malloc(1, int(ret.size))
	}

	//arg typs ptr
	ats := new(unsafe.Pointer)
	if len(args) == 0 {
		_ats := mem.Malloc(1, 8)
		(*ats) = _ats
		mem.PushAt(_ats, 0, unsafe.Pointer(TypeVoid.toC()))
	} else {
		_ats := mem.Malloc(len(args), 8)
		(*ats) = _ats
		for i := range args {
			mem.PushAt(_ats, i, unsafe.Pointer(args[i].toC()))
		}
	}

	_ret := TypeVoid
	if ret != nil {
		_ret = ret
	}

	st := C.ffi_prep_cif(cif.cif, abi.toC(),
		C.uint(len(args)), _ret.toC(), (**C.ffi_type)(*ats))

	err := Status(st).Error()
	if err != nil {
		defer cif.Free()
		return nil, err
	}

	return cif, nil
}
func (cif *Cif) Free() {
	if cif.retVal != nil {
		mem.Free(cif.retVal)
	}
	if cif.cif.arg_types != nil {
		mem.Free(unsafe.Pointer(cif.cif.arg_types))
	}
	if cif.cif != nil {
		mem.Free(unsafe.Pointer(cif.cif))
	}
}
func (cif *Cif) Call(fn unsafe.Pointer, argAddr ...any) unsafe.Pointer {
	argc := len(argAddr)
	argv := new(unsafe.Pointer)
	if argc == 0 {
		_argv := mem.Malloc(1, 8)
		(*argv) = _argv
		defer mem.Free(_argv)
	} else {
		_argv := mem.Malloc(argc, 8)
		(*argv) = _argv
		defer mem.Free(_argv)

		for i := range argAddr {
			if argAddr[i] == nil {
				mem.PushAt(_argv, i, nil)
				continue
			}
			mem.PushAt(_argv, i, reflect.ValueOf(argAddr[i]).UnsafePointer())
		}
	}
	C.ffi_call(cif.cif, (*[0]byte)(fn), cif.retVal, (*unsafe.Pointer)(*argv))
	return cif.retVal
}

type Closure struct {
	cif *Cif

	closure *C.ffi_closure
	fnptr   unsafe.Pointer
	data    unsafe.Pointer
}

type ClosureConf struct {
	Abi  Abi
	Args []*Type
	Ret  *Type
}
type ClosureData struct {
	Args     []unsafe.Pointer
	Ret      unsafe.Pointer
	UserData []any
}
type closureUserData struct {
	fn       func(*ClosureData)
	argc     int
	userData *[]any
}

//export closure_caller
func closure_caller(cif *C.ffi_cif, ret, args, userData unsafe.Pointer) {
	data := (*closureUserData)(userData)
	input := new(ClosureData)
	input.Args = *(*[]unsafe.Pointer)(mem.Slice(args, data.argc))
	input.Ret = ret
	if data.userData != nil {
		input.UserData = *data.userData
	}
	data.fn(input)
}
func NewClosure(conf ClosureConf, fn func(*ClosureData), userData ...any) *Closure {
	var err error
	cls := new(Closure)
	cls.cif, err = NewCif(conf.Abi, conf.Ret, conf.Args...)
	if err != nil {
		panic(err)
	}
	cls.fnptr = mem.Malloc(1, 8)
	cls.closure = (*C.ffi_closure)(C.ffi_closure_alloc(
		C.ulonglong(mem.Sizeof(C.ffi_closure{})), (*unsafe.Pointer)(cls.fnptr)))

	cls.data = mem.Malloc(1, mem.Sizeof(closureUserData{}))
	data := (*closureUserData)(cls.data)
	data.fn = fn
	data.argc = len(conf.Args)
	if userData != nil {
		data.userData = &userData
	} else {
		data.userData = nil
	}

	C.ffi_prep_closure_loc(cls.closure, cls.cif.cif,
		(*[0]byte)(C.closure_caller), cls.data, mem.Pop(cls.fnptr))
	return cls
}
func (cls *Closure) Call(args ...any) unsafe.Pointer {
	return cls.cif.Call(mem.Pop(cls.fnptr), args...)
}
func (cls *Closure) Free() {
	if cls.cif != nil {
		cls.cif.Free()
	}
	if cls.fnptr != nil {
		mem.Free(cls.fnptr)
	}
	if cls.data != nil {
		data := (*closureUserData)(cls.data)
		data.fn = nil
		data.userData = nil
		mem.Free(cls.data)
	}
	if cls.closure != nil {
		C.ffi_closure_free(unsafe.Pointer(cls.closure))
	}
}
func (cls *Closure) Cfn() unsafe.Pointer {
	return mem.Pop(cls.fnptr)
}
