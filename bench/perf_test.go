package main

import (
	"fmt"
	"testing"

	"github.com/hk-32/utils/out"
)

// fmt causes an allocation per element for slices inside of arrays & slice... very expensive!
// var value = [...]any{100, 200, 400, "Hello", "World", []int{100, 200, 300, 400, 500, 600}}
var value = []any{100, 200, 400, "Hello", "World", []int{100, 200, 300, 400, 500, 600}}

//var value any = []string{"Hello", "World"}

//var value any = []int{100, 200, 300, 400, 500, 600}

//var value any = 25

/* var v1 = []int{50, 100, 50, 100}
var value = []any{"Hello", v1} */

// 275 - 279

func BenchmarkStandard(b *testing.B) {
	b.ReportAllocs()
	for i := 1; i < b.N; i++ {
		//_ = fmt.Sprintf("%v is a %v who is %v years old!", "Hassan", "male", 19)
		_ = fmt.Sprint(value)
	}
}

func BenchmarkCustom(b *testing.B) {
	b.ReportAllocs()
	for i := 1; i < b.N; i++ {
		//_ = Sprintf("%v is a %v who is %v years old!", "Hassan", "male", 19)
		_ = out.Stringify(value)
	}
}
