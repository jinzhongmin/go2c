package main

/*
#include <stdio.h>
int add(int a, int b){
	return a+b;
}
*/
import "C"
import (
	"fmt"

	"github.com/jinzhongmin/go2c/pkg/ffi"
)

func main() {
	cif, _ := ffi.NewCif(ffi.AbiDefault, //abi
		ffi.TypeInt32,                //return type
		ffi.TypeInt32, ffi.TypeInt32) //arg type

	a := int32(100)
	b := int32(200)

	c := cif.Call(C.add, &a, &b)
	fmt.Println(*(*int)(c))

	cif.Free()
}
