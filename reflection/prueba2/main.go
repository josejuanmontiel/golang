package main

import (
	"reflect"
    "strconv"
	"net/http"
	"fmt"
	"runtime"
	"errors"
)


type Interface interface{}

type Struct struct{}

type Type2 struct {
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
	typeReg.Set(new(Type2))
	runtime.GC()
}



func main() {

	route("/",handler1)
	route2("/handler2",handler2)
	
	FuncAnalyse(handler2)

	// one way is to have a value of the type you want already
	a := 1
	// reflect.New works kind of like the built-in function new
	// We'll get a reflected pointer to a new int value
	intPtr := reflect.New(reflect.TypeOf(a))
	// Just to prove it
	b := intPtr.Elem().Interface().(int)
	// Prints 0
	fmt.Println(b)

	// We can also use reflect.New without having a value of the type
	var nilInt *int
	intType := reflect.TypeOf(nilInt).Elem()
	intPtr2 := reflect.New(intType)
	// Same as above
	c := intPtr2.Elem().Interface().(int)
	// Prints 0 again
	fmt.Println(c)

}

func route(path string, function func(w http.ResponseWriter, r *http.Request)) {
	fmt.Println("\n--->Run route")
	function(nil,nil)
}

func route2(path string, function func(w http.ResponseWriter, r *http.Request, s Type2)) {
	fmt.Println("\n--->Run route 2")


	abc, err := typeReg.Get("Type2")
	if err != nil {
		fmt.Println(err)
		return
	}

	function(nil,nil, abc)
}

func handler1(w http.ResponseWriter, r *http.Request){
	fmt.Println("\n--->Handler 1")
}

func handler2(w http.ResponseWriter, r *http.Request, body Type2){
	fmt.Println("\n--->Handler 2")
}

/*
var typeRegistry = make(map[string]reflect.Type)

func init() {
    myTypes := []interface{}{MyString{}}
    for _, v := range myTypes {
        // typeRegistry["MyString"] = reflect.TypeOf(MyString{})
        typeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
    }
}
*/

func FuncAnalyse(m interface{}) {

	//Reflection type of the underlying data of the interface
	x := reflect.TypeOf(m)

	numIn := x.NumIn() //Count inbound parameters
	numOut := x.NumOut() //Count outbounding parameters

	fmt.Println("Method:", x.String()) 
	fmt.Println("Variadic:", x.IsVariadic()) // Used (<type> ...) ?
	fmt.Println("Package:", x.PkgPath())


	for i := 0; i < numIn; i++ {
		inV := x.In(i)
		in_Kind := inV.Kind() //func
		fmt.Printf("\nParameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\n-----------",in_Kind,inV.Name())

		if "Type2"==inV.Name() {
			var nilType2 *Type2
			intType2 := reflect.TypeOf(nilType2).Elem()
			intPtr2 := reflect.New(intType2)
			// Same as above
			abc := intPtr2.Elem().Interface().(Type2)
			handler2(nil,nil, abc)
		}

	}

	for o := 0; o < numOut; o++ {
		returnV := x.Out(0)
		return_Kind := returnV.Kind()
		fmt.Printf("\nParameter OUT: "+strconv.Itoa(o)+"\nKind: %v\nName: %v\n",return_Kind,returnV.Name())
	}

}
