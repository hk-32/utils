package out

import (
	"strconv"
	"unsafe"

	"github.com/hk-32/utils/rt"
)

// Stringify formats value using the default format and returns the resulting string.
func Stringify(value any) string {
	if value == nil {
		return "<nil>"
	}

	// Maybe it can describe itself
	if v, isStringer := value.(Stringer); isStringer {
		return v.String()
	} else if v2, isErr := value.(error); isErr {
		return v2.Error()
	}

	iface := (*rt.AnyInternal)(unsafe.Pointer(&value))

	switch rt.KindOf(value) {
	case rt.String:
		v := *(*string)(iface.Data)
		return v
	case rt.Bool:
		v := *(*bool)(iface.Data)
		if v {
			return "true"
		}
		return "false"

	// Integers
	case rt.Int:
		v := *(*int)(iface.Data)
		return strconv.Itoa(v)
	case rt.Int8:
		v := *(*int8)(iface.Data)
		return strconv.Itoa(int(v))
	case rt.Int16:
		v := *(*int16)(iface.Data)
		return strconv.Itoa(int(v))
	case rt.Int32:
		v := *(*int32)(iface.Data)
		return strconv.Itoa(int(v))
	case rt.Int64:
		v := *(*int64)(iface.Data)
		return strconv.FormatInt(v, 10)

	// Floats
	case rt.Float32:
		v := *(*float32)(iface.Data)
		sp := make([]byte, 0, 23) // -> stack allocation
		return string(strconv.AppendFloat(sp, float64(v), 'g', -1, 32))
		//return strconv.FormatFloat(float64(v), 'g', -1, 32) // Causes an allocation
	case rt.Float64:
		v := *(*float64)(iface.Data)
		sp := make([]byte, 0, 23) // -> stack allocation
		return string(strconv.AppendFloat(sp, v, 'g', -1, 64))
		//return strconv.FormatFloat(v, 'g', -1, 64) // Causes an allocation

	// Unsigned Integers
	case rt.Uint:
		v := *(*uint)(iface.Data)
		return strconv.FormatUint(uint64(v), 10)
	case rt.Uint8:
		v := *(*uint8)(iface.Data)
		return strconv.FormatUint(uint64(v), 10)
	case rt.Uint16:
		v := *(*uint16)(iface.Data)
		return strconv.FormatUint(uint64(v), 10)
	case rt.Uint32:
		v := *(*uint32)(iface.Data)
		return strconv.FormatUint(uint64(v), 10)
	case rt.Uint64:
		v := *(*uint64)(iface.Data)
		return strconv.FormatUint(v, 10)
	case rt.Uintptr:
		v := *(*uintptr)(iface.Data)
		return strconv.FormatUint(uint64(v), 10)

	case rt.Complex64:
		v := *(*complex64)(iface.Data)
		return strconv.FormatComplex(complex128(v), 'f', -1, 64)
	case rt.Complex128:
		v := *(*complex128)(iface.Data)
		return strconv.FormatComplex(v, 'f', -1, 128)

	case rt.Array, rt.Slice:
		buffer := make([]byte, 0, space(value))
		guessedSpace := cap(buffer)
		buffer = stringify_into(value, buffer)
		if guessedSpace != len(buffer) {
			panic("There is a bug in the space counting process of slices/arrays!")
		}
		return *(*string)(unsafe.Pointer(&buffer))

	case rt.Pointer:
		iface.Type = (*rt.PointerType)(iface.Type).PointedToType
		// basically only go 1 pointer deep... and even then, only if the pointed to value isn't an any itself
		if rt.KindOf(value) != rt.Pointer && rt.KindOf(value) != rt.Interface {
			return "&" + Stringify(value)
		}
		// just return the address instead
		// Max uint64 = 18446744073709551615 = FFFFFFFFFFFFFFFF = len() = 16
		var buffer = make([]byte, 0, 16+2) // -> stack allocation
		buffer = strconv.AppendUint(append(buffer, '0', 'x'), uint64(uintptr(iface.Data)), 16)
		return string(buffer)

	case rt.UnsafePointer:
		// just return the address of unsafe pointers
		var buffer = make([]byte, 0, 16+2) // -> stack allocation
		buffer = strconv.AppendUint(append(buffer, '0', 'x'), uint64(uintptr(iface.Data)), 16)
		return string(buffer)

	case rt.Interface:
		inner_face := *(*any)(iface.Data)
		return Stringify(inner_face)

	case rt.Struct:
		return "<struct>"
	}
	return "<" + rt.KindOf(value).String() + ">"
}

// stringify_into is entirely used within this package, it can only append to a pre-made buffer.
// so a user call to 'Stringify' will allocate the right length buffer and then further calls will be to this function.
// This saves us from an allocation per stringification like Stringify does
func stringify_into(value any, buffer []byte) []byte {
	if value == nil {
		return append(buffer, "<nil>"...)
	}

	// Maybe it can describe itself
	if v, isStringer := value.(Stringer); isStringer {
		return append(buffer, v.String()...)
	} else if v2, isErr := value.(error); isErr {
		return append(buffer, v2.Error()...)
	}

	iface := (*rt.AnyInternal)(unsafe.Pointer(&value))

	switch rt.KindOf(value) {
	case rt.String:
		v := *(*string)(iface.Data)
		return append(buffer, v...)
	case rt.Bool:
		v := *(*bool)(iface.Data)
		if v {
			return append(buffer, "true"...)
		}
		return append(buffer, "false"...)

	// Integers
	case rt.Int:
		v := *(*int)(iface.Data)
		return strconv.AppendInt(buffer, int64(v), 10)
	case rt.Int8:
		v := *(*int8)(iface.Data)
		return strconv.AppendInt(buffer, int64(v), 10)
	case rt.Int16:
		v := *(*int16)(iface.Data)
		return strconv.AppendInt(buffer, int64(v), 10)
	case rt.Int32:
		v := *(*int32)(iface.Data)
		return strconv.AppendInt(buffer, int64(v), 10)
	case rt.Int64:
		v := *(*int64)(iface.Data)
		return strconv.AppendInt(buffer, v, 10)

	// Floats
	case rt.Float32:
		v := *(*float32)(iface.Data)
		return strconv.AppendFloat(buffer, float64(v), 'g', -1, 32)
	case rt.Float64:
		v := *(*float64)(iface.Data)
		return strconv.AppendFloat(buffer, v, 'g', -1, 64)

	// Unsigned Integers
	case rt.Uint:
		v := *(*uint)(iface.Data)
		return strconv.AppendUint(buffer, uint64(v), 10)
	case rt.Uint8:
		v := *(*uint8)(iface.Data)
		return strconv.AppendUint(buffer, uint64(v), 10)
	case rt.Uint16:
		v := *(*uint16)(iface.Data)
		return strconv.AppendUint(buffer, uint64(v), 10)
	case rt.Uint32:
		v := *(*uint32)(iface.Data)
		return strconv.AppendUint(buffer, uint64(v), 10)
	case rt.Uint64:
		v := *(*uint64)(iface.Data)
		return strconv.AppendUint(buffer, v, 10)
	case rt.Uintptr:
		v := *(*uintptr)(iface.Data)
		return strconv.AppendUint(buffer, uint64(v), 10)

	case rt.Complex64:
		v := *(*complex64)(iface.Data)
		return append(buffer, strconv.FormatComplex(complex128(v), 'f', -1, 64)...)
	case rt.Complex128:
		v := *(*complex128)(iface.Data)
		return append(buffer, strconv.FormatComplex(v, 'f', -1, 128)...)

	case rt.Array:
		array := (*rt.ArrayType)(iface.Type)

		// Construct a holder for the array element values
		container_internal := rt.AnyInternal{Type: array.ElemsType}
		container := (*any)(unsafe.Pointer(&container_internal))

		elems_kind := rt.KindOf((*container))
		size := rt.SizeOfAny((*container))

		buffer = append(buffer, '[')
		for i := 0; i < int(array.Len); i++ {
			container_internal.Data = unsafe.Add(iface.Data, size*i)

			// Don't double wrap values in 'any' if it itself is one
			if elems_kind == rt.Interface {
				container = (*any)(unsafe.Pointer(container_internal.Data))
			}

			buffer = stringify_into(*container, buffer)

			// Add spacing
			if i != int(array.Len)-1 {
				buffer = append(buffer, ' ')
			}
		}
		return append(buffer, ']')

	/* case rt.Array:
	// Quickly convert the array into a slice and fallthrough to slice handling
	array := (*rt.ArrayType)(iface.Type)
	iface.Data = unsafe.Pointer(&rt.SliceData{iface.Data, int(array.Len), int(array.Len)})
	iface.Type = array.SliceType
	fallthrough */

	case rt.Slice:
		slice := (*rt.SliceData)(iface.Data)

		// Construct a holder for the slice element values
		container_internal := rt.AnyInternal{Type: (*rt.SliceType)(iface.Type).ElemsType}
		container := (*any)(unsafe.Pointer(&container_internal))

		elems_kind := rt.KindOf((*container))
		size := rt.SizeOfAny((*container))

		buffer = append(buffer, '[')
		for i := 0; i < slice.Len; i++ {
			container_internal.Data = unsafe.Add(slice.Data, size*i)

			// Don't double wrap values in 'any' if it itself is one
			if elems_kind == rt.Interface {
				container = (*any)(unsafe.Pointer(container_internal.Data))
			}

			buffer = stringify_into(*container, buffer)

			// Add spacing
			if i != slice.Len-1 {
				buffer = append(buffer, ' ')
			}
		}
		return append(buffer, ']')

	case rt.Pointer:
		iface.Type = (*rt.PointerType)(unsafe.Pointer(iface.Type)).PointedToType
		// basically only go 1 pointer deep... and even then, only if the pointed to value isn't an any itself
		if rt.KindOf(value) != rt.Pointer && rt.KindOf(value) != rt.Interface {
			return stringify_into(value, append(buffer, '&'))
		}
		// just print the address instead
		return strconv.AppendUint(append(buffer, '0', 'x'), uint64(uintptr(iface.Data)), 16)

	case rt.UnsafePointer:
		// just return the address of unsafe pointers
		return strconv.AppendUint(append(buffer, '0', 'x'), uint64(uintptr(iface.Data)), 16)

	case rt.Interface:
		inner_face := *(*any)(iface.Data)
		return stringify_into(inner_face, buffer)

	case rt.Struct:
		return append(buffer, "<struct>"...)
	}
	return append(append(append(buffer, '<'), rt.KindOf(value).String()...), '>')

}
