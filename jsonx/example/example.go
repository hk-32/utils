package main

import (
	"fmt"

	"github.com/hk-32/utils/jsonx"
	"github.com/hk-32/utils/out"
)

var schema = jsonx.Object{
	"name":    "Bob",
	"age":     16,
	"friends": []any{"Hassan", "Sufian", "Arild", "Osama"},
}

type fruit string
type amount int

var schema2 = map[fruit]amount{
	"Apples":  2,
	"Oranges": 6,
	"Bananas": 3,
}

func main() {
	value := []any{"Hassan", "Sufian", "Arild", "Osama", true, false, 2.5, nil}
	fmt.Println(jsonx.Stringify(value))

	out.Println(jsonx.Stringify(schema))
	fmt.Println(schema)
	fmt.Println(jsonx.GenericMapToJSON(schema2))

	/* file, err := ioutil.ReadFile("./input.json")
	if err != nil {
		panic(err)
	}

	//jsonx.ByteDigitsToNumber = jsonx.ByteDigitsToInt
	structure, err := jsonx.Decode(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(structure.(jsonx.Object)) */

	// Don't know types... receiving arbitrary json
	/*if name, isSet := structure.(jsonx.Object)["name"]; isSet {
		if name, isStr := name.(string); isStr {
			fmt.Println("NAME IS A STRING :", name)
		}
	}*/

	/*if jsonx.Match(schema, structure) {
		// Use type assertions without fear
		person := structure.(jsonx.Object)

		fmt.Printf("Person's name is %v\n", person["name"])
		fmt.Printf("Person's age is %v\n", person["age"])
	} else {
		fmt.Println("Schema didn't match!")
	}*/
}
