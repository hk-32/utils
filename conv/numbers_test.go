package conv

import (
	"strconv"
	"testing"
)

var buffer1 = make([]byte, 0, 20)

func BenchmarkStandard(b *testing.B) {
	b.ReportAllocs()
	for i := 1; i < b.N; i++ {
		_ = len(strconv.AppendInt(buffer1, 19, 10))
	}
}

func BenchmarkCustom(b *testing.B) {
	b.ReportAllocs()
	for i := 1; i < b.N; i++ {
		_ = i64bytes(19)
	}
}

func i64bytes(x int64) int {
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

/* func BenchmarkCustom(b *testing.B) {
	b.ReportAllocs()
	for i := 1; i < b.N; i++ {
		_ = Itoa(92233720368547)
	}
} */
