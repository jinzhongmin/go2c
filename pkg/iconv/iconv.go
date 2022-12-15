package iconv

/*
#cgo LDFLAGS: -liconv -lcharset
#include <stdlib.h>
#include <iconv.h>
#include <libcharset.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/jinzhongmin/go2c/pkg/mem"
)

func LocaleCharset() string {
	lc := C.locale_charset()
	return C.GoString(lc)
}

type Iconvt struct {
	c       C.iconv_t
	inptr   unsafe.Pointer
	insize  unsafe.Pointer
	outptr  unsafe.Pointer
	outsize unsafe.Pointer
}

func Open(toCode, fromCode string) (*Iconvt, error) {
	tc := C.CString(toCode)
	defer C.free(unsafe.Pointer(tc))
	fc := C.CString(fromCode)
	defer C.free(unsafe.Pointer(fc))

	c := C.iconv_open(tc, fc)
	if c == nil {
		return nil, errors.New("create iconvt err : toCode or fromCode error")
	}
	it := new(Iconvt)
	it.c = c
	it.inptr = mem.Malloc(1, 8)
	it.insize = mem.Malloc(1, 8)
	it.outptr = mem.Malloc(1, 8)
	it.outsize = mem.Malloc(1, 8)
	return it, nil
}
func (it *Iconvt) IconvChars(data unsafe.Pointer, l int) ([]byte, error) {
	bufSize := l

	inbuf := data
	*(*C.ulonglong)(it.insize) = C.ulonglong(bufSize)

	outbuf := mem.Malloc(bufSize, 1)
	defer mem.Free(outbuf)
	*(*C.ulonglong)(it.outsize) = C.ulonglong(bufSize)
	outslice := *(*[]byte)(mem.Slice(outbuf, bufSize))

	mem.PushTo(it.inptr, inbuf)
	mem.PushTo(it.outptr, outbuf)

	output := make([]byte, 0)
	for {
		inputs := *(*C.ulonglong)(it.insize)

		r := int(C.iconv(it.c, (**C.char)(it.inptr), (*C.ulonglong)(it.insize),
			(**C.char)(it.outptr), (*C.ulonglong)(it.outsize)))
		if r == -1 && inputs > 0 && inputs == *(*C.ulonglong)(it.insize) {
			return nil, errors.New("in Iconv looks like the codes don't match")
		}

		output = append(output, outslice[:C.ulonglong(bufSize)-(*(*C.ulonglong)(it.outsize))]...)
		if (*(*C.ulonglong)(it.insize)) == 0 {
			break
		}

		cvts := bufSize - int(*(*C.ulonglong)(it.insize))
		mem.PushTo(it.inptr, unsafe.Pointer(uintptr(inbuf)+uintptr(cvts)))
		mem.PushTo(it.outptr, outbuf)
		*(*C.ulonglong)(it.outsize) = C.ulonglong(bufSize)
	}
	return output, nil
}
func (it *Iconvt) IconvBytes(data []byte) ([]byte, error) {
	bufSize := len(data)

	inbuf := unsafe.Pointer(&data[0])
	*(*C.ulonglong)(it.insize) = C.ulonglong(bufSize)

	outbuf := mem.Malloc(bufSize, 1)
	defer mem.Free(outbuf)
	*(*C.ulonglong)(it.outsize) = C.ulonglong(bufSize)
	outslice := *(*[]byte)(mem.Slice(outbuf, bufSize))

	mem.PushTo(it.inptr, inbuf)
	mem.PushTo(it.outptr, outbuf)

	output := make([]byte, 0)
	for {
		inputs := *(*C.ulonglong)(it.insize)

		r := int(C.iconv(it.c, (**C.char)(it.inptr), (*C.ulonglong)(it.insize),
			(**C.char)(it.outptr), (*C.ulonglong)(it.outsize)))
		if r == -1 && inputs > 0 && inputs == *(*C.ulonglong)(it.insize) {
			return nil, errors.New("in Iconv looks like the codes don't match")
		}

		output = append(output, outslice[:C.ulonglong(bufSize)-(*(*C.ulonglong)(it.outsize))]...)
		if (*(*C.ulonglong)(it.insize)) == 0 {
			break
		}

		cvts := bufSize - int(*(*C.ulonglong)(it.insize))
		mem.PushTo(it.inptr, unsafe.Pointer(uintptr(inbuf)+uintptr(cvts)))
		mem.PushTo(it.outptr, outbuf)
		*(*C.ulonglong)(it.outsize) = C.ulonglong(bufSize)
	}
	return output, nil
}
func (it *Iconvt) IconvStr(str string) ([]byte, error) {
	type _str struct {
		str unsafe.Pointer
		len int
	}
	gostr := (*_str)(unsafe.Pointer(&str))

	bufSize := gostr.len
	inbuf := gostr.str

	*(*C.ulonglong)(it.insize) = C.ulonglong(bufSize)

	outbuf := mem.Malloc(bufSize, 1)
	defer mem.Free(outbuf)
	*(*C.ulonglong)(it.outsize) = C.ulonglong(bufSize)
	outslice := *(*[]byte)(mem.Slice(outbuf, bufSize))

	mem.PushTo(it.inptr, inbuf)
	mem.PushTo(it.outptr, outbuf)

	output := make([]byte, 0)
	for {
		inputs := *(*C.ulonglong)(it.insize)

		r := int(C.iconv(it.c, (**C.char)(it.inptr), (*C.ulonglong)(it.insize),
			(**C.char)(it.outptr), (*C.ulonglong)(it.outsize)))
		if r == -1 && inputs > 0 && inputs == *(*C.ulonglong)(it.insize) {
			return nil, errors.New("in Iconv looks like the codes don't match")
		}

		output = append(output, outslice[:C.ulonglong(bufSize)-(*(*C.ulonglong)(it.outsize))]...)
		if (*(*C.ulonglong)(it.insize)) == 0 {
			break
		}

		cvts := bufSize - int(*(*C.ulonglong)(it.insize))
		mem.PushTo(it.inptr, unsafe.Pointer(uintptr(inbuf)+uintptr(cvts)))
		mem.PushTo(it.outptr, outbuf)
		*(*C.ulonglong)(it.outsize) = C.ulonglong(bufSize)
	}
	return output, nil
}
func (it *Iconvt) Close() {
	mem.Free(it.inptr)
	mem.Free(it.insize)
	mem.Free(it.outptr)
	mem.Free(it.outsize)
	C.iconv_close(it.c)
}

func Simply(src string, toCode, fromCode string) (string, error) {
	it, err := Open(toCode, fromCode)
	if err != nil {
		return "", err
	}
	defer it.Close()

	r, err := it.IconvStr(src)
	if err != nil {
		return "", err
	}
	return string(r), nil
}
