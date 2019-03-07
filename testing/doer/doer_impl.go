package doer

import (
	"errors"
	"fmt"
)

type doerImpl struct {
}

func NewDoerService() Doer {
	return &doerImpl{}
}

func (d *doerImpl) DoSomething(numero int, cadena string) error {
	fmt.Println("Calling DoSomething method")
	fmt.Println(cadena)
	return errors.New("error")
}
