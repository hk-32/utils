package main

import "github.com/hk-32/utils/out"

/*
NOTE:
	fmt.Println vs out.Println = 357kb difference (unstripped) Go 1.20.2
*/

func main() {
	//os.Stdout.WriteString("Hello World") 1483 KB
	out.Println("Hello World") // -> 1567 KB
	//println("Hello World") // -> 1215 KB
	//fmt.Println("Hello World") // -> 1924 KB

	// Prints the same output as if it was 'fmt.Println'
	//out.Println("Hello World", true, 100, vector2{420, 69}, nil, []int8{1, 10, 20}, []any{50, 20, "Bye", "World", []string{"Yo"}})
	//fmt.Println("Hello World", true, 100, vector2{420, 69}, nil, []int8{1, 10, 20}, []any{50, 20, "Bye", "World", []string{"Yo"}})

	/*out.Printf("%v is a %v who is %v years old!\n", "Hassan", "male", 19)

	/* x := []string{"Hello", "World"}
	var value any = []any{100, 200, 300, complex128(45), 500, 600, -9223372036854775807, unsafe.Pointer(&x)} */

	/* var value any = []any{100, 200, 400, "Hello", "World", &vector2{25, 50}, []int{100, 200, 300, 400, 500, 600}}

	fmt.Println(fmt.Sprint(value))
	out.Println(out.Stringify(value))

	/*
		name := in.ReadLine("Enter your name: ")
		age := in.ReadLine("Enter your age: ")
		out.Printf("Hello %v, you're %v years old!\n", name, age)
	*/

	/*
		age, err := strconv.ParseInt(in.Input("Please enter your age: "), 10, 64)
		if err != nil {
			println("Invalid age!")
		}
		println(age)
	*/
}

/* type vector2 struct {
	x int
	y int
}

// struct types require a method 'String' that can describe the structure
func (v vector2) String() string {
	return out.Group(v.x, v.y)
} */

/* func (v vector2) Format(buffer out.Formatter) {
	buffer.AppendGroup(v.x, v.y)
} */
