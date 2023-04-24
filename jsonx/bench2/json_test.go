package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/hk-32/utils/jsonx"
)

const fileN = "./input.json"

func BenchmarkCustom(b *testing.B) {
	file, err := ioutil.ReadFile(fileN)
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		CustomDecode(file)
	}
}

func CustomDecode(in []byte) interface{} {
	structure, err := jsonx.Decode(in)
	if err != nil {
		panic(err)
	}
	return structure
}

func BenchmarkStandard(b *testing.B) {
	file, err := ioutil.ReadFile(fileN)
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		StandardDecode(file)
	}
}

func StandardDecode(in []byte) interface{} {
	var structure interface{}
	err := json.Unmarshal(in, &structure)
	if err != nil {
		panic(err)
	}
	return structure
}
