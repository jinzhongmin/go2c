package main

import "C"
import (
	"fmt"

	"github.com/jinzhongmin/go2c/pkg/dlfcn"
	"github.com/jinzhongmin/go2c/pkg/ffi"
)

//export test
func test() {
	fmt.Println("test")
}

//For presentation purposes,
//the program loads itself into memory and
//reads the symbol from it. And make the call using ffi.Cif.
func main() {
	dl, _ := dlfcn.DlOpen("./loaddl.exe", dlfcn.RTLDLazy)
	fn, _ := dl.Symbol("test")

	cif, _ := ffi.NewCif(ffi.AbiDefault, nil)
	cif.Call(fn)

	cif.Free()
	dl.Close()
}
