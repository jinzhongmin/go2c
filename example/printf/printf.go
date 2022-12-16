package main

//#include <stdio.h>
//void* func_printf(){ return printf; } //get printf pointer
import "C"
import (
	"unsafe"

	"github.com/jinzhongmin/go2c/pkg/ffi"
	"github.com/jinzhongmin/go2c/pkg/mem"
)

func main() {
	cif, _ := ffi.NewCif(ffi.AbiDefault, nil,
		ffi.TypePointer, ffi.TypePointer, ffi.TypeInt32)
	defer cif.Free()

	format := C.CString("hello, im %s, %d years old.\r\n")
	lilei := C.CString("lilei")
	defer mem.Free(unsafe.Pointer(format))
	defer mem.Free(unsafe.Pointer(lilei))

	var old int32 = 18
	cif.Call(C.func_printf(), &format, &lilei, &old)
}
