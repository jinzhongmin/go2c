package ffi

/*
#include <ffi.h>
*/
import "C"
import (
	"errors"
)

type Type C.ffi_type

func (typ *Type) toC() *C.ffi_type {
	return (*C.ffi_type)(typ)
}
func (typ *Type) String() string {
	switch typ {
	case TypeVoid:
		return "void"
	case TypePointer:
		return "pointer"
	case TypeUint8:
		return "uint8"
	case TypeInt8:
		return "int8"
	case TypeUint16:
		return "uint16"
	case TypeInt16:
		return "int16"
	case TypeUint32:
		return "uint32"
	case TypeInt32:
		return "int32"
	case TypeUint64:
		return "uint64"
	case TypeInt64:
		return "int64"
	}
	return "unknow"
}

var (
	TypeVoid    *Type = (*Type)(&C.ffi_type_void)
	TypePointer *Type = (*Type)(&C.ffi_type_pointer)
	TypeUint8   *Type = (*Type)(&C.ffi_type_uint8)
	TypeInt8    *Type = (*Type)(&C.ffi_type_sint8)
	TypeUint16  *Type = (*Type)(&C.ffi_type_uint16)
	TypeInt16   *Type = (*Type)(&C.ffi_type_sint16)
	TypeUint32  *Type = (*Type)(&C.ffi_type_uint32)
	TypeInt32   *Type = (*Type)(&C.ffi_type_sint32)
	TypeUint64  *Type = (*Type)(&C.ffi_type_uint64)
	TypeInt64   *Type = (*Type)(&C.ffi_type_sint64)

	TypeFloat             *Type = (*Type)(&C.ffi_type_float)
	TypeDouble            *Type = (*Type)(&C.ffi_type_double)
	TypeLongDouble        *Type = (*Type)(&C.ffi_type_longdouble)
	TypeComplexFloat      *Type = (*Type)(&C.ffi_type_complex_float)
	TypeComplexdouble     *Type = (*Type)(&C.ffi_type_complex_double)
	TypeComplexLongdouble *Type = (*Type)(&C.ffi_type_complex_longdouble)
)

type Abi C.ffi_abi

func (a Abi) toC() C.ffi_abi {
	return C.ffi_abi(a)
}

const (
	AbiFirstABI Abi = 0
	AbiSysv     Abi = 1
	AbiThiscall Abi = 3
	AbiFastcall Abi = 4
	AbiStdcall  Abi = 5
	AbiPascal   Abi = 6
	AbiRegister Abi = 7
	AbiMsCdecl  Abi = 8
	AbiDefault  Abi = AbiSysv
)

type Status C.ffi_status

func (st Status) Error() error {
	switch C.ffi_status(st) {
	case C.FFI_OK:
		return nil
	case C.FFI_BAD_TYPEDEF:
		return errors.New("bad typedef")
	case C.FFI_BAD_ABI:
		return errors.New("bad abi")
	case C.FFI_BAD_ARGTYPE:
		return errors.New("bad argtype")
	}
	return nil
}
