package main

import (
	"reflect"
    "strconv"
	"fmt"
	"runtime"
	"errors"
)

type ExampleType struct {
	attr1 string
}

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
	typeReg.Set(new(ExampleType))
	runtime.GC()
}

type Handler struct{}

func main() {
	h := new(Handler)

	h.route("/handler",method1)
}

func (h Handler) route(path string, function func(w string, r string, s ExampleType)) {
	fmt.Println("\n--->Run route")

	abc, err := typeReg.Get(GetNameLastParam(function))
	if err != nil {
		fmt.Println(err)
		return
	}

	method := reflect.ValueOf(h).MethodByName("method1")
	in := make([]reflect.Value, method.Type().NumIn())
	
	in[0] = reflect.ValueOf("1")
	in[1] = reflect.ValueOf("2")
	in[2] = reflect.ValueOf(abc)

	response := method.Call(in)
	fmt.Println("--------->")
	fmt.Println(response)
	fmt.Println("<---------")
}

func method1(w string, r string, body ExampleType){
	fmt.Println("\n--->Handler")
}

func GetNameLastParam(m interface{}) (string) {

	//Reflection type of the underlying data of the interface
	x := reflect.TypeOf(m)

	numIn := x.NumIn() //Count inbound parameters

	inV := x.In(numIn-1)
	in_Kind := inV.Kind() //func

	fmt.Printf("\nParameter IN: "+strconv.Itoa(numIn-1)+"\nKind: %v\nName: %v\n-----------",in_Kind,inV.Name())

	return inV.Name()
}
