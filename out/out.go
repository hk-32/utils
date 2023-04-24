package out

import (
	"os"
	"unsafe"

	"github.com/hk-32/utils/rt"
)

/* func init() {
	file, err := os.Create("output")
	if err != nil {
		panic(err)
	}
	os.Stdout = file
} */

// stringer - Useful for custom types or structs formatting.
type Stringer interface {
	String() string
}

func Printf(format string, a ...any) (int, error) {
	return os.Stdout.WriteString(Sprintf(format, a...))
}

// Println formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended. It returns the number of bytes written and any write error encountered.
func Println(a ...any) (int, error) {
	length := 0
	for i, value := range a {
		length += space(value)

		// Account for whitespace if not last argument
		if i != len(a)-1 {
			length++
		}
	}

	buffer := make([]byte, length+1)
	cursor := 0

	for i, value := range a {
		cursor += copy(buffer[cursor:], Stringify(value))

		// Add spaces
		if i != len(a)-1 {
			buffer[cursor] = ' '
			cursor++
		}
	}

	buffer[cursor] = '\n'
	return os.Stdout.Write(buffer)
}

/* func Println(a ...any) (int, error) {
	var bytes []byte

	for i, value := range a {
		bytes = append(bytes, Stringify(value)...)

		// Add a spacing
		if i != len(a)-1 {
			bytes = append(bytes, ' ')
		}
	}
	return os.Stdout.Write(append(bytes, '\n'))
} */

// Print Formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands. It returns the number of bytes written and any write error encountered.
func Print(a ...any) (int, error) {
	var bytes []byte

	for i, value := range a {
		bytes = append(bytes, Stringify(value)...)

		// Add a spacing
		if i != len(a)-1 {
			bytes = append(bytes, ' ')
		}
	}
	return os.Stdout.Write(bytes)
}

// Sprintf can be used for string interpolation, use '%v' for slot markers.
func Sprintf(format string, a ...any) string {
	// calculate the extra lenght that the arguments will add
	added_length := 0
	for _, v := range a {
		added_length += space(v)
	}
	slots_len := len(a) * 2

	buffer := make([]byte, 0, len(format)+added_length-slots_len)

	slots_count := 0
	for i := 0; i < len(format); i++ {
		// if not '%' then its just a normal character
		if format[i] != '%' {
			buffer = append(buffer, format[i])
			continue
		}

		// if this is the last character then append and break
		if i == len(format)-1 {
			buffer = append(buffer, '%')
			break
		}
		// proceed to the verb
		i++

		// if this slot doesn't have an argument left to input
		if slots_count >= len(a) {
			buffer = append(buffer, '%', format[i])
			continue
		}

		switch format[i] {
		case 'v':
			buffer = append(buffer, Stringify(a[slots_count])...)
			slots_count++
		case 'T':
			buffer = append(buffer, rt.KindOf(a[slots_count]).String()...)
			slots_count++
		default:
			buffer = append(buffer, '%', format[i])
		}
	}

	return *(*string)(unsafe.Pointer(&buffer))
}

// Group can be used to group a bunch of values and surround them with '{ }'.
// Useful for 'String' methods on structs.
func Group(args ...any) string {
	bytes := make([]byte, 1, 12)
	bytes[0] = '{'

	for i, value := range args {
		// Add value
		bytes = append(bytes, Stringify(value)...)

		// Add a spacing
		if i != len(args)-1 {
			bytes = append(bytes, ' ')
		}
	}

	return string(append(bytes, '}'))
}
