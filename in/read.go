package in

import (
	"os"

	"github.com/hk-32/utils/out"
)

var buffer = make([]byte, 32)

// ReadLine reads input until an 'enter' key.
func ReadLine(prompt ...interface{}) string {
	if prompt != nil {
		out.Print(prompt...)
	}

	var bytes []byte

	for {
		r, err := os.Stdin.Read(buffer)
		if err != nil {
			break
		}

		// input is only returned to us when the user presses the 'enter' key.
		// a '\r\n' or just a '\n' indicates the end
		if buffer[r-1] == '\n' {
			if buffer[r-2] == '\r' {
				bytes = append(bytes, buffer[:r-2]...)
			} else {
				bytes = append(bytes, buffer[:r-1]...)
			}
			break
		}

		// Not the end! Keep reading
		bytes = append(bytes, buffer[:r]...)
	}
	return string(bytes)
}
