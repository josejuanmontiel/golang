package main

import (
	"github.com/josejuanmontiel/golang/packaging/dep"
)

func main() {
	newDepImpl := dep.NewDepService()
	_ = newDepImpl.DoSomething(1, "hola")
}
