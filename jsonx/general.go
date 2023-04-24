package jsonx

import (
	"errors"
	"unsafe"
)

// ErrDefault - Default error for Decode()
var ErrDefault = errors.New("JSON could not be parsed")

// ErrInvalidNumber - Number wrongly formatted
var ErrInvalidNumber = errors.New("invalid number format")

// ErrNumberOutOfRange - Number exceeds the range of its holder
var ErrNumberOutOfRange = errors.New("number exceeds the range of an 'int64'")

// state - Internal structure for keeping track of state
type state struct {
	src []byte // The whole input
	pos int    // Current position in source
	len int    // Lenght of source

	refer2src bool // If using 'FlagReferSrc'
	str2num   func(bytes []byte) (interface{}, error)
}

// Object represents a JSON Object
type Object map[string]any

// String encodes the object as JSON and returns it
func (obj Object) String() string {
	buffer := make([]byte, 0, space(obj))
	buffer = stringify_into(obj, buffer)
	return *(*string)(unsafe.Pointer(&buffer))
}

// Array - Represents a JSON Array
type Array = []any

// Identifiers
const (
	// EOS - Used internally to signify end of stream
	iEOS byte = 0x03 // 0x03 = End of Text, 0x04 = End of Transmission
	// Star - Used internally to support jsonc: JSON Comments
	iStar byte = '*'
	// Slash - Used internally to support jsonc: JSON Comments
	iFSlash byte = '/'

	// Dot - Used internally for reading floats
	iDot byte = '.'
	// Quotation - Used internally for reading strings
	iQuotation byte = '"'

	// Hyphen - Syntax literal negative
	iHyphen byte = '-'

	// Comma - Syntax literal comma
	iComma byte = ','
	// Colon - Syntax literal colon
	iColon byte = ':'

	// LeftBrace - Syntax literal to start an object
	iLeftBrace byte = '{'
	// RightBrace - Syntax literal to end an object
	iRightBrace byte = '}'

	// LeftBracket - Syntax literal to start a list
	iLeftBracket byte = '['
	// RightBracket - Syntax literal to end a list
	iRightBracket byte = ']'
)

// get returns the value of a given index in source
func (state *state) get(n int) byte {
	if n < state.len {
		return state.src[n]
	}
	return iEOS
}

// peek returns the first non-space without advancing the position
func (state *state) peek() byte {
	for i := state.pos; i < state.len; i++ {
		if isSpace(state.src[i]) {
			return state.src[i]
		}
	}
	return iEOS
}

// swallow returns the first non-space and advances the position
func (state *state) swallow() byte {
	for state.pos+1 < state.len {
		state.pos++
		if !isSpace(state.src[state.pos]) {
			return state.src[state.pos]
		}
	}
	return iEOS
}

// isDigit checks if the byte is a digit from 0 - 9
func isDigit(r byte) bool {
	// return '0' <= r && r <= '9'
	return r >= '0' && r <= '9'
}

// isSpace checks if the byte is empty space
func isSpace(r byte) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

/* // 1 with 18 zeroes is the max limit of int64
var pow10tab = [...]int64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
}

// pow10tabf - Stores the pre-computed values upto 10^(31) or 1 with 31 zeros
var pow10tabf = [...]float64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22, 1e23, 1e24, 1e25, 1e26, 1e27, 1e28, 1e29,
	1e30, 1e31,
}

// pow10 - Raises 10 by n
func pow10(n int) int64 {
	if n > 31 {
		n = 31
	}
	return pow10tab[n]
}

func pow10f(n int) float64 {
	//if n >= 0 && n <= 31 { // n >= 0 && n <= 308
	// return pow10postab32[uint(n)/32] * pow10tab[uint(n)%32]
	if n > 31 {
		n = 31
	}
	return pow10tabf[n]
	//}
	//panic("out of range number encountered")
} */
