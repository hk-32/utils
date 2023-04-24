package jsonx

import (
	"strconv"
	"unsafe"

	"github.com/hk-32/utils/out"
	"github.com/hk-32/utils/rt"
	"golang.org/x/exp/constraints"
)

type Jsonify interface {
	Jsonify() string
}

type SimpleValue interface {
	constraints.Integer | constraints.Float | ~string | ~bool | Object
}

// GenericMapToJSON can take any single layer Go map, provided that the key is of an underlying type string
func GenericMapToJSON[K ~string, V SimpleValue](obj map[K]V) string {
	// Calculate space
	lenght := 2
	i := 0
	for k, v := range obj {
		lenght += space(k) + 2 + space(v)
		// Account for spaces & commas between values
		if i != len(obj)-1 {
			lenght += 2
		}
		i++
	}

	buffer := make([]byte, 0, lenght)
	buffer = append(buffer, '{')
	i = 0

	for k, v := range obj {
		// Key
		buffer = append(append(append(append(buffer, '"'), k...), '"'), ':', ' ')
		// Value
		buffer = stringify_into(v, buffer)

		// Account for spaces and commas between values
		if i != len(obj)-1 {
			buffer = append(buffer, ',', ' ')
		}
		i++
	}
	buffer = append(buffer, '}')
	if lenght != len(buffer) {
		panic("There is a bug in the space counting process of array/slice!")
	}
	return *(*string)(unsafe.Pointer(&buffer))
}

// Stringify converts a jsonx.Object or any Go slice/array into JSON
func Stringify(value any) string {
	if value == nil {
		return "<nil>"
	}

	switch rt.KindOf(value) {
	case rt.Array, rt.Slice, rt.Map:
		buffer := make([]byte, 0, space(value))
		guessed := cap(buffer)
		buffer = stringify_into(value, buffer)
		if guessed != len(buffer) {
			panic(out.Sprintf("There is a bug in the space calculating process of %Ts!", value))
		}
		return *(*string)(unsafe.Pointer(&buffer))
	}
	panic("Stringify(): Unsuported value of type <" + rt.KindOf(value).String() + ">")
}

// stringify_into is entirely used within this package, it can only append to a pre-made buffer.
// so a user call to 'Stringify' will allocate the right length buffer and then further calls will be to this function.
// This saves us from an allocation per stringification like Stringify does
func stringify_into(value any, buffer []byte) []byte {
	if value == nil {
		return append(buffer, "null"...)
	}

	// Maybe it can Jsonify itself
	if value, isJsonible := value.(Jsonify); isJsonible {
		return append(buffer, value.Jsonify()...)
	}

	iface := (*rt.AnyInternal)(unsafe.Pointer(&value))

	switch rt.KindOf(value) {
	case rt.String:
		v := *(*string)(iface.Data)
		return append(append(append(buffer, '"'), v...), '"')
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

	case rt.Array:
		array := (*rt.ArrayType)(iface.Type)

		// Construct a holder for the array element values
		container_internal := &rt.AnyInternal{Type: array.ElemsType, Data: iface.Data}
		container := (*any)(unsafe.Pointer(container_internal))

		elems_kind := rt.KindOf((*container))
		size := rt.SizeOf((*container))

		buffer = append(buffer, '[')
		for i := 0; i < int(array.Len); i++ {
			container_internal.Data = unsafe.Add(iface.Data, size*i)

			// Don't double wrap values in the type any if it already is an any
			if elems_kind == rt.Interface {
				container = (*any)(unsafe.Pointer(container_internal.Data))
			}

			buffer = stringify_into(*container, buffer)

			// Account for spaces and commas between values
			if i != int(array.Len)-1 {
				buffer = append(buffer, ',', ' ')
			}
		}
		return append(buffer, ']')

	case rt.Slice:
		slice := (*rt.SliceData)(iface.Data)

		// Construct a holder for the slice element values
		container_internal := &rt.AnyInternal{Type: (*rt.SliceType)(iface.Type).ElemsType, Data: slice.Data}
		container := (*any)(unsafe.Pointer(container_internal))

		elems_kind := rt.KindOf((*container))
		size := rt.SizeOf((*container))

		buffer = append(buffer, '[')
		for i := 0; i < slice.Len; i++ {
			container_internal.Data = unsafe.Add(slice.Data, size*i)

			// Don't double wrap values in the type any if it already is an any
			if elems_kind == rt.Interface {
				container = (*any)(unsafe.Pointer(container_internal.Data))
			}

			buffer = stringify_into(*container, buffer)

			// Account for spaces and commas between values
			if i != slice.Len-1 {
				buffer = append(buffer, ',', ' ')
			}
		}
		return append(buffer, ']')

	case rt.Map:
		obj, isObj := value.(Object)
		if !isObj {
			panic("Stringify(): Given map type isn't jsonx.Object!")
		}

		buffer = append(buffer, '{')
		i := 0
		for k, v := range obj {
			// Key
			buffer = append(append(append(append(buffer, '"'), k...), '"'), ':', ' ')
			// Value
			buffer = stringify_into(v, buffer)

			// Account for spaces and commas between values
			if i != len(obj)-1 {
				buffer = append(buffer, ',', ' ')
			}
			i++
		}
		return append(buffer, '}')

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
