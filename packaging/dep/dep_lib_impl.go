package dep

import (
	"errors"
	"fmt"
)

type depImpl struct {
}

func NewDepService() DepLib {
	return &depImpl{}
}

func (d *depImpl) DoSomething(numero int, cadena string) error {
	fmt.Println("Calling DoSomething method")
	fmt.Println(cadena)
	return errors.New("error")
}
