package jsonx

import (
	"unsafe"
)

// Flag - Options for the Decoder
type Flag int8

// Flags
const (
	// FlagRefer2Src - Refer to source for strings instead of allocating copies
	FlagRefer2Src Flag = iota

	// FlagCopy - Can be used in conjunction with 'FlagRefer2Src' to use a copy of source input,
	// this makes the resulting json independant of the input slice '[]byte'
	FlagCopy

	// FlagNumbersAsF64 - (Recommended) Return all numbers as 64 bit floating points 'float64' just like in javascript
	FlagNumbersAsF64

	// FlagNumbersAsI64 - Return all numbers as 64 bit integers 'int64'. Results in decimal loss aka 2.5 becomes 2
	FlagNumbersAsI64
)

// Decode - Parses and decodes the input. Flags are for specific cases!
func Decode(input []byte, flags ...Flag) (interface{}, error) {
	// Create a default decoder
	dec := state{src: input, pos: -1, len: len(input)}
	dec.str2num = stof

	// Parse flags
	for _, v := range flags {
		if v == FlagCopy {
			dec.src = make([]byte, len(input))
			copy(dec.src, input)
		} else if v == FlagRefer2Src {
			dec.refer2src = true
		} else if v == FlagNumbersAsI64 {
			dec.str2num = stoi
		} else if v == FlagNumbersAsF64 {
			dec.str2num = stof
		}
	}

	// Compose the structure
	out, err := dec.compose()
	// Just to make sure that out is nothing except for nil on failure
	if err != nil {
		return nil, err
	}
	return out, nil
}

// compose - Returns Objects, Arrays, Strings and Numbers with err nil, rest as actual char but err
func (state *state) compose() (interface{}, error) {
	char := state.swallow()

	switch char {
	case iQuotation:
		// Get string length
		var strLen int
	COUNT:
		char = state.get(state.pos + strLen + 1)
		if char == iEOS {
			// Unfinished strings are illegal, so error out
			return nil, ErrDefault
		} else if char != iQuotation {
			strLen++
			goto COUNT
		}
		// Now update state's position to the trailing quotation
		state.pos += strLen + 1

		// Now separate string chars
		if state.refer2src {
			str := state.src[state.pos-strLen : state.pos]
			return *(*string)(unsafe.Pointer(&str)), nil
		}
		// Default: return a copy; string() results in an allocation as it duplicates the source
		return string(state.src[state.pos-strLen : state.pos]), nil

	case iLeftBracket:
		var List Array

	PARSE_ITEM:
		element, err := state.compose()
		if err == nil {
			List = append(List, element)

			// Next should be an iComma or an iRightBracket
			switch state.swallow() {
			case iComma: // There is more to come...
				goto PARSE_ITEM
			case iRightBracket: // The end has been reached.
				return List, nil
			default: // Error; Unexpected Token: Expected a Comma or a RightBracket
				return nil, ErrDefault
			}
		} else if element == iRightBracket && len(List) == 0 {
			// Array had nothing...
			return List, nil
		}
		// Error ... element can be iEOS or something that does not make sense here
		return nil, err

	case iLeftBrace:
		var Map = make(Object)

	PARSE_PAIR:
		// Parse key
		char = state.swallow()
		if char == iQuotation {
			// Get string length
			var strLen int
		COUNT_KEY_LEN:
			char = state.get(state.pos + strLen + 1)
			if char == iEOS {
				// Unfinished strings are illegal, so error out
				return iEOS, ErrDefault
			} else if char != iQuotation {
				strLen++
				goto COUNT_KEY_LEN
			}
			// Now update state's position to the trailing quotation
			state.pos += strLen + 1

			var key []byte
			// Now separate string chars
			if state.refer2src {
				key = state.src[state.pos-strLen : state.pos]
			} else {
				key = make([]byte, strLen)
				copy(key, state.src[state.pos-strLen:state.pos])
			}

			// Swallow a colon
			if state.swallow() == iColon {
				// Compose value
				value, err := state.compose()
				if err == nil {
					Map[*(*string)(unsafe.Pointer(&key))] = value

					// Swallow next... should be an iComma or an iRightBrace
					char = state.swallow()
					if char == iComma {
						// There is more to come...
						goto PARSE_PAIR
					} else if char == iRightBrace {
						// The end has been reached.
						return Map, nil
					}
					// Unexpected Token: Expected a ',' or a '}'
					return nil, ErrDefault
				}
				return value, err
			}
			// Error Unexpected Token: Expected a Colon
			return nil, ErrDefault
		} else if char == iRightBrace && len(Map) == 0 {
			// Object had nothing...
			return Map, nil
		}
		// Error ... next can be iEOS or something that does not make sense here
		return nil, ErrDefault

	case 't':
		if state.pos+3 < state.len {
			if state.src[state.pos+1] == 'r' &&
				state.src[state.pos+2] == 'u' &&
				state.src[state.pos+3] == 'e' {
				state.pos += 3
				return true, nil
			}
		}
	case 'f':
		if state.pos+4 < state.len {
			if state.src[state.pos+1] == 'a' &&
				state.src[state.pos+2] == 'l' &&
				state.src[state.pos+3] == 's' &&
				state.src[state.pos+4] == 'e' {
				state.pos += 4
				return false, nil
			}
		}
	case 'n':
		if state.pos+3 < state.len {
			if state.src[state.pos+1] == 'u' &&
				state.src[state.pos+2] == 'l' &&
				state.src[state.pos+3] == 'l' {
				state.pos += 3
				return nil, nil
			}
		}
	}

	if isDigit(char) || char == iHyphen {
		var numLen int = 1
		var hasDecimals bool = false

		for char = state.get(state.pos + 1); char != iEOS; char = state.get(state.pos + numLen) {
			if isDigit(char) {
				numLen++
				continue
			} else if char == iDot {
				if !hasDecimals {
					numLen++
					hasDecimals = true
					continue
				}
				// Apparently more than one dot were found
				return nil, ErrDefault
			}

			// Done with parsing number
			num := state.src[state.pos : state.pos+numLen]
			// Now update state's position
			// -1 because state.pos starts with the position of the first digit so account for that
			state.pos += numLen - 1

			if state.str2num != nil {
				return state.str2num(num)
			}
			return Number{hasDecimals, num}, nil
		}
		return iEOS, ErrDefault
	}
	return char, ErrDefault
}

// Specifically written for compose's number parser
func stoi(str []byte) (any, error) {
	// Formats the digits according to base 10
	var power int64 = 1
	var number int64

	// Reverse loop & add digit adjusted for power to total
	for i := len(str) - 1; i >= 0; i-- {
		char := str[i]

		// It's a digit
		if char >= '0' && char <= '9' {
			number += (int64(char-'0') * power)
			if number >= 0 {
				// Ok
				power *= 10
				continue
			} else if (i == 1 && str[0] == '-') && number == -9223372036854775808 {
				// Bit of a special case for -9223372036854775808 as its actually allowed
				return number, nil
			}

			// Number exceeds the maximum allowed
			return 0, ErrNumberOutOfRange
		} else if char == '-' && i == 0 {
			return -number, nil
		}

		// Anything else is obviously wrongly fomatted
		return 0, ErrInvalidNumber
	}
	return number, nil
}

// Specifically written for compose's number parser
func stof(str []byte) (any, error) {
	// Formats the digits according to base 10
	var power float64 = 1
	var number float64

	// Reverse loop & add digit adjusted for power to total
	for i := len(str) - 1; i >= 0; i-- {
		char := str[i]

		if (char >= '0') && (char <= '9') {
			// It's a digit
			number += (float64(char - '0')) * power
			power *= 10
			continue
		} else if char == '-' && i == 0 {
			return -number, nil
		} else if char == '.' {
			// Divide number by it's upper power & reset power
			number /= power
			power = 1
			continue
		}

		// Anything else is obviously wrongly fomatted
		return 0, ErrInvalidNumber
	}
	return number, nil
}

/* // Specifically written for compose's number parser
func stof(chars []char) interface{} {
	// Formats the digits according to base 10
	var exponent int
	var number float64

	// Reverse loop to start from the right & add digit adjusted for power to total
	for i := len(chars) - 1; i >= 0; i-- {
		char := chars[i]

		if isDigit(char) {
			number += (float64(char - '0')) * pow10f(exponent)
			exponent++
		} else if char == iHyphen && i == 0 {
			return -number
		} else if char == iDot {
			// Divide number by its upper length & reset exponent
			number /= pow10f(exponent)
			exponent = 0
		}
	}
	return number
}

// Specifically written for compose's number parser
func stoi(chars []char) interface{} {
	// Formats the digits according to base 10
	var exponent int
	var number int64

	// Reverse loop to start from the right & add digits adjusted for power to total
	for i := len(chars) - 1; i >= 0; i-- {
		char := chars[i]

		if isDigit(char) {
			number += (int64(char - '0')) * pow10(exponent)
			exponent++
		} else if char == '-' && i == 0 {
			return -number
		}
	}
	return number
} */
