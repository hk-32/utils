package out

import (
	"strconv"
	"unsafe"

	"github.com/hk-32/utils/rt"
	"golang.org/x/exp/constraints"
)

// space calculates the space a value will take as a string in bytes
func space(value any) int {
	if value == nil {
		return 5
	}

	// Maybe it can describe itself
	if v, isStringer := value.(Stringer); isStringer {
		return len(v.String())
	} else if v2, isErr := value.(error); isErr {
		return len(v2.Error())
	}

	iface := (*rt.AnyInternal)(unsafe.Pointer(&value))

	switch rt.KindOf(value) {
	case rt.String:
		v := *(*string)(iface.Data)
		return len(v)
	case rt.Bool:
		v := *(*bool)(iface.Data)
		if v {
			return len("true")
		}
		return len("false")

	// Integers
	case rt.Int:
		v := *(*int)(iface.Data)
		return numSpace(v)
	case rt.Int8:
		v := *(*int8)(iface.Data)
		return numSpace(v)
	case rt.Int16:
		v := *(*int16)(iface.Data)
		return numSpace(v)
	case rt.Int32:
		v := *(*int32)(iface.Data)
		return numSpace(v)
	case rt.Int64:
		v := *(*int64)(iface.Data)
		return numSpace(v)

	// Floats
	case rt.Float32:
		sp := make([]byte, 0, 23) // -> stack allocation
		v := *(*float32)(iface.Data)
		return len(strconv.AppendFloat(sp, float64(v), 'g', -1, 32))
	case rt.Float64:
		sp := make([]byte, 0, 23) // -> stack allocation
		v := *(*float64)(iface.Data)
		return len(strconv.AppendFloat(sp, v, 'g', -1, 64))

	// Unsigned Integers
	case rt.Uint:
		v := *(*uint)(iface.Data)
		return numSpace(v)
	case rt.Uint8:
		v := *(*uint8)(iface.Data)
		return numSpace(v)
	case rt.Uint16:
		v := *(*uint16)(iface.Data)
		return numSpace(v)
	case rt.Uint32:
		v := *(*uint32)(iface.Data)
		return numSpace(v)
	case rt.Uint64:
		v := *(*uint64)(iface.Data)
		return numSpace(v)
	case rt.Uintptr:
		v := *(*uintptr)(iface.Data)
		return numSpace(v)

	case rt.Complex64:
		v := *(*complex64)(iface.Data)
		// expensive but we don't have AppendComplex
		return len(strconv.FormatComplex(complex128(v), 'f', -1, 64))
	case rt.Complex128:
		v := *(*complex128)(iface.Data)
		return len(strconv.FormatComplex(v, 'f', -1, 128))

	case rt.Array:
		array := (*rt.ArrayType)(iface.Type)

		// Construct a holder for the array element values
		container_internal := &rt.AnyInternal{Type: array.ElemsType, Data: iface.Data}
		container := (*any)(unsafe.Pointer(container_internal))

		elems_kind := rt.KindOf(*container)
		size := rt.SizeOfAny((*container))

		added_length := 0
		for i := 0; i < int(array.Len); i++ {
			container_internal.Data = unsafe.Add(iface.Data, size*i)

			// Don't double wrap values in the type any if it already is an any
			if elems_kind == rt.Interface {
				container = (*any)(unsafe.Pointer(container_internal.Data))
			}

			added_length += space(*container)

			// Add spacing
			if i != int(array.Len)-1 {
				added_length++
			}
		}
		return added_length + 2

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

		elems_kind := rt.KindOf(*container)
		size := rt.SizeOfAny(*container)

		added_length := 0
		for i := 0; i < slice.Len; i++ {
			container_internal.Data = unsafe.Add(slice.Data, size*i)

			// Don't double wrap values in 'any' if it itself is one
			if elems_kind == rt.Interface {
				container = (*any)(unsafe.Pointer(container_internal.Data))
			}

			added_length += space(*container)
			// Account for spaces between values
			if i < slice.Len-1 {
				added_length++
			}
		}
		return added_length + 2

	case rt.Pointer:
		iface.Type = (*rt.PointerType)(unsafe.Pointer(iface.Type)).PointedToType
		// basically only go 1 pointer deep... and even then, only if the pointed to value isn't an any itself
		if rt.KindOf(value) != rt.Pointer && rt.KindOf(value) != rt.Interface {
			return space(value) + 1
		}
		// just return the address lenght instead
		sp := make([]byte, 0, 18)
		return len(strconv.AppendUint(sp, uint64(uintptr(iface.Data)), 16)) + 2

	case rt.UnsafePointer:
		sp := make([]byte, 0, 18)
		return len(strconv.AppendUint(sp, uint64(uintptr(iface.Data)), 16)) + 2

	case rt.Interface:
		inner_face := *(*any)(iface.Data)
		return space(inner_face)

	case rt.Struct:
		return len("<struct>")
	}
	return int(len(rt.KindOf(value).String())) + 2
}

func numSpace[T constraints.Integer](x T) int {
	if x == 0 {
		return 1
	}

	count := 0
	if x < 0 {
		// The '-' sign will also take a byte
		count++
	}
	for x > 0 || x < 0 {
		x = x / 10
		count++
	}

	return count
}

/* func i64bytes(x int64) int {
	if x == 0 {
		return 1
	}

	count := 0
	if x < 0 {
		// The '-' sign will also take a byte
		count++
	}
	for x > 0 || x < 0 {
		x = x / 10
		count++
	}

	return count
}

func u64bytes(x uint64) int {
	if x == 0 {
		return 1
	}

	count := 0
	for x > 0 {
		x = x / 10
		count++
	}

	return count
}
*/
