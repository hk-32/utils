package conv

import "errors"

var ErrInvalidSyntax = errors.New("invalid Syntax")
var ErrOutOfRange = errors.New("number exceeds the maximum range of an 'int64'")

// StrToInt - Converts a base 10 number from string to an int64
func StrToInt(str string) (int64, error) {
	// Formats the digits according to base 10
	var power int64 = 1
	var number int64

	// Reverse loop & add digit adjusted for power to total
	for i := len(str) - 1; i >= 0; i-- {
		char := str[i]

		if char >= '0' && char <= '9' {
			// It's a digit
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
			return 0, ErrOutOfRange
		} else if char == '-' && i == 0 {
			return -number, nil
		}

		// Anything else is obviously wrongly fomatted
		return 0, ErrInvalidSyntax
	}
	return number, nil
}

// StrToFloat - Converts a base 10 number from string to a float64
func StrToFloat(str string) (float64, error) {
	// Formats the digits according to base 10
	var power float64 = 1
	var number float64
	var hasDecimals bool

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
		} else if char == '.' && !hasDecimals {
			// Divide number by it's upper power & reset power
			number /= power
			power = 1
			hasDecimals = true
			continue
		}

		// Anything else is obviously wrongly fomatted
		return 0, ErrInvalidSyntax
	}
	return number, nil
}

// Simplest version without any error checking; just for reference
func stoi(bytes string) int64 {
	// Formats the digits according to base 10
	var exponent int64 = 1
	var number int64

	// Reverse loop to start from the right & add digits adjusted for power to total
	for i := len(bytes) - 1; i >= 0; i-- {
		char := bytes[i]

		if char >= '0' && char <= '9' {
			number += (int64(char - '0')) * exponent
			exponent *= 10
		} else if char == '-' && i == 0 {
			return -number
		}
	}
	return number
}
