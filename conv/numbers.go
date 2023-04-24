package conv

const nSmalls = 100

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

var buffer [20]byte

// small returns the string for an i with 0 <= i < nSmalls.
func small(i int) string {
	if i < 10 {
		return digits[i : i+1]
	}
	return smallsString[i*2 : i*2+2]
}

/* func _itoa(num int64) string {
	//buffer := make([]byte, 20)
	var isNegative bool = num < 0

	// We can only handle numbers in positive
	if isNegative {
		num = -num
	}

	// Process individual digits
	i := 19
	for num != 0 {
		buffer[i] = byte(num%10) + '0'
		num /= 10
		i--
	}

	if isNegative {
		buffer[i] = '-'
		i--
	}

	slice := buffer[i+1:]
	//return string(buffer[i+1:])
	return *(*string)(unsafe.Pointer(&slice))
}

type Integer interface {
	int8 | int16 | int32 | int64 | int
}

func Itoa[T Integer](x T) string {
	if 0 <= x && x < nSmalls {
		return small(int(x))
	}
	return _itoa(int64(x))
} */

func AppendItoa(num int64, buffer []byte) []byte {
	if 0 <= num && num < nSmalls {
		return append(buffer, small(int(num))...)
	}

	if num == 0 {
		return append(buffer, '0')
	}
	var temp [20]byte
	var isNegative bool = num < 0

	// We can only handle numbers in positive
	if isNegative {
		num = -num
	}

	// Process individual digits
	i := 19
	for num != 0 {
		temp[i] = byte(num%10) + '0'
		num /= 10
		i--
	}

	if isNegative {
		temp[i] = '-'
		i--
	}

	return append(buffer, temp[i+1:]...)
}
