package main

import (
	"fmt"
	"reflect"
)

type uuid string

func typeOf(v interface{}) string {
	t := reflect.TypeOf(v)
	return t.Name()
}

func main() {
	var x uuid = "test"
	fmt.Println(typeOf(x))
}
