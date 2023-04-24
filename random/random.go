package random

import (
	"math/rand"
)

/* func main() {
	var chance0, chance1 int

	for i := 0; i < 1000; i++ {
		switch Pick(0, 1) {
		case 0:
			chance0++
		case 1:
			chance1++
		}
	}

	fmt.Printf("Chances of 0 are %vx\n", chance0)
	fmt.Printf("Chances of 1 are %vx\n", chance1)
} */

/* func init() {
	x := rand.NewRand(NewSource(seed))
} */

// takes a variadic number of arguments and returns one of them randomly
func Pick[T any](items ...T) T {
	return items[rand.Intn(len(items))]
}
