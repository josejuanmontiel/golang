package main

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

type X struct{}

type TypeRegister map[string]reflect.Type

func (t TypeRegister) Set(i interface{}) {
	//name string, typ reflect.Type
	t[reflect.TypeOf(i).Name()] = reflect.TypeOf(i)
}

func (t TypeRegister) Get(name string) (interface{}, error) {
	if typ, ok := t[name]; ok {
		return reflect.New(typ).Elem().Interface(), nil
	}
	return nil, errors.New("no one")
}

var typeReg = make(TypeRegister)

func init() {
	typeReg.Set(new(X))
	runtime.GC()
}

func main() {
	fmt.Println("Create new *X element")
	t := new(X)
	fmt.Printf("Name of new *X element type is '%v'\n", reflect.TypeOf(t).Name())
	y, err := typeReg.Get(reflect.TypeOf(t).Name())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Name of new y element type is '%v'\n", reflect.TypeOf(y).Name())
	fmt.Println(typeReg)
	fmt.Printf("Tell me what the hell is wrong with it: %v\n", reflect.TypeOf(y).Name())
	fmt.Printf("Tell me what the hell is wrong with it: %T\n", y)
	fmt.Printf("Tell me what the hell is wrong with it: %v\n", reflect.TypeOf(y).String())
}