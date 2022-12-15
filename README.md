# go2c

go2c 旨在解决使用CGO过程中遇到的一些问题，算是一个功能合集，对一些C库进行了绑定。

go2c is designed to solve some of the problems encountered in using CGO. It is a collection of features that bind some C libraries.

程序是在windows下编写的，考虑到windows用户可能没有安装相关的C库，3rdparty 目录下附带了编译好的二进制文件。

The program was written under windows, and since windows users may not have the relevant C library installed, the compiled binaries are included in the 3rdparty directory.

## mem 

github.com/jinzhongmin/go2c/pkg/mem

mem 是对内存操作的库。特别是指针的嵌套容易搞混。通过Push相关方法，将一个指针的地址压入另一个指针。

mem is a library for manipulating memory. In particular, the nesting of Pointers can be confusing. The address of one pointer is pushed into another by the Push correlation method.

```C
void *p;
void **p;
p[0] = p;
```

在mem中使用下面的代码达到前面C代码的功能。

Use the following code in mem to achieve the functionality of the previous C code.

```golang
var p  unsafe.Pointer
var pp unsafe.Pointer
mem.PushTo(pp, p)

//你也可以这样写，但是这样写pp可能被GC释放掉，这样的指针也不能传到C代码中，可能会有问题。
//You can also write this, but writing pp like this may be released by GC, and such Pointers cannot be passed to C code, which may cause problems.

var p  unsafe.Pointer
pp := new(unsafe.Pointer)
*(*unsafe.Pointer)(pp) = p

```

## ffi 

github.com/jinzhongmin/go2c/pkg/ffi

ffi 是为了解决怎么向C中传递回调函数的问题。它是对libffi的golang绑定。

ffi is designed to solve the problem of how to pass a callback function to C. It is a golang binding to libffi.

libffi 的使用及相关许可请访问 https://github.com/libffi/libffi

The use of libffi and related license please visit https://github.com/libffi/libffi

通过Cif可以构建一个调用c的模板，然后在go中调用c的函数。可以参考example 中的 add，Cif调用的参数和返回参数必须是指针。

Using Cif you can build a template that calls c and then calls c's function in go. As you can see from add in example, the parameters and return parameters of the Cif call must be Pointers.

通过Closure可以将回调函数传递到C中去。可以参考example 中的 callback。

Closure allows you to pass callback functions into C. Refer to the callback in the example.

## dlfcn

github.com/jinzhongmin/go2c/pkg/dlfcn

dlfcn 是一个动态加载C共享库函数的库，返回的Symbol指针可以通过ffi的Cif进行调用。他是对dlfcn的绑定。

dlfcn is a library that dynamically loads C shared library functions and returns a Symbol pointer that can be called with the Cif of ffi.It is a binding to dlfcn.

附带的 windows 二进制库编译自 https://github.com/dlfcn-win32/dlfcn-win32 , 相关使用及许可请访问获取。

The included windows binary library is compiled from https://github.com/dlfcn-win32/dlfcn-win32 and can be used and licensed here.

通过 dlfcn 加载的函数需要使用ffi的Cif进行调用。

Functions loaded via dlfcn need to be called using the Cif of ffi.

## iconv

github.com/jinzhongmin/go2c/pkg/iconv

iconv 是对 iconv c库的绑定，提供了字符集的转换功能。

iconv is a binding to the iconv c library that provides character set conversion.

附带的 windows 二进制库编译自 https://www.gnu.org/software/libiconv/, 相关使用及许可请访问获取。

The included windows binary library is compiled from https://www.gnu.org/software/libiconv/dlfcn-win32 and can be used and licensed here.

