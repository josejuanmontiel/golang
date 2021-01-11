package main

import (
	"fmt"
	"reflect"
)

type Aaa struct {
	a string
}

type Bbb struct {
	b int
}

type Handler struct{}

func (h Handler) GET(a Aaa, b Bbb, ptr *Aaa) string {
	return "OK" + a.a + " ptr:" + ptr.a
}

func main() {
	handler := new(Handler)

	objects := make(map[reflect.Type]interface{})
	objects[reflect.TypeOf(Aaa{})] = Aaa{"jkljkL"}
	objects[reflect.TypeOf(new(Aaa))] = &Aaa{"pointer!"}
	objects[reflect.TypeOf(Bbb{})] = Bbb{}

	//in := make([]reflect.Value, 0)
	method := reflect.ValueOf(handler).MethodByName("GET")

	fmt.Println(method)

	in := make([]reflect.Value, method.Type().NumIn())

	fmt.Println("method type num in:", method.Type().NumIn())
	for i := 0; i < method.Type().NumIn(); i++ {
		t := method.Type().In(i)
		object := objects[t]
		fmt.Println(i, "->", object)
		in[i] = reflect.ValueOf(object)
	}

	fmt.Println("method type num out:", method.Type().NumOut())

	response := method.Call(in)
	fmt.Println(response)
}