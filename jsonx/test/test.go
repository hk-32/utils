package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hk-32/utils/jsonx"
)

var schema = jsonx.Object{
	"name": "Bob",
	"age":  float64(16),
}

func main() {
	file, err := ioutil.ReadFile("./input.json")
	if err != nil {
		panic(err)
	}

	structure, err := jsonx.Decode(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(structure)
}
