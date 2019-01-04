package main

import (
	"fmt"
	"reflect"
)

func checkFunc(f interface{}) {
	// 1. check type
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		fmt.Printf("please check f is not nil")
	}

	// 2. check in input params
	ft := reflect.TypeOf(f)

	if ft.NumIn() != 3 {
		fmt.Printf("the f should have 3 params")
	}

	if ft.In(0).Kind() != reflect.Int {
		fmt.Printf("the f first param should be int\n")
	}

	if !ft.In(2).Implements(reflect.TypeOf((*Interface)(nil)).Elem()) {
		fmt.Printf("the 3th param should impl Interface\n")
	}

	// 3. check output params
	if ft.NumOut() != 2 {
		fmt.Printf("the number of return should 2")
	}

	if ft.Out(0).Kind() != reflect.Int {
		fmt.Printf("the fisrt return should int")
	}

	if ft.Out(1).Kind() != reflect.String {
		fmt.Printf("the second return should string")
	}
}

type Interface interface {
	Hw()
}

type hw struct {
}

func (*hw) Hw() {

}

func add(x, y int, hw *hw) (int, string) {
	return 0, ""
}

func main() {
	checkFunc(add)
}
