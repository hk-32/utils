package jsonx

import "strconv"

// Number represents a JSON number
type Number struct {
	hasDecimals bool
	bytes       []byte
}

// String returns the Go string version of it
func (n Number) String() string {
	return string(n.bytes)
}

// HasDecimals - Returns if the number has decimal digits
func (n Number) HasDecimals() bool {
	return n.hasDecimals
}

// AsFloat64 - Returns it as a float64
func (n Number) AsFloat64() (float64, error) {
	return strconv.ParseFloat(string(n.bytes), 64)
}

// AsInt64 - Returns it as a int64
func (n Number) AsInt64() (int64, error) {
	return strconv.ParseInt(string(n.bytes), 10, 64)
}
