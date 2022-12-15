package main

/*
typedef int (*add)(int a, int b);

int dofn(void *p, int a, int b){
	add fn = (add)p;
	return fn(a, b);
}
*/
import "C"
import (
	"fmt"

	"github.com/jinzhongmin/go2c/pkg/ffi"
)

func main() {
	cls := ffi.NewClosure(
		//The template for configuring Closure is the same as int (*add)(int a, int b)
		ffi.ClosureConf{
			Abi: ffi.AbiDefault,
			Ret: ffi.TypeInt32,
			Args: []*ffi.Type{
				ffi.TypeInt32, ffi.TypeInt32,
			},
			//Body of function
		}, func(cd *ffi.ClosureData) {
			a := *(*int)(cd.Args[0]) //get the first arg
			b := *(*int)(cd.Args[1]) //get the second arg

			*(*int)(cd.Ret) = a + b //set the return of the function
		})

	//You can call this function directly
	a := 100
	b := 200
	c := cls.Call(&a, &b)
	fmt.Println("call by go, result is", *(*int)(c))

	//You can also pass a function into C and call it through C
	r := C.dofn(cls.Cfn(), 100, 200)
	fmt.Println("call by c callback, result is", r)

	//Remember to release
	cls.Free()

}
