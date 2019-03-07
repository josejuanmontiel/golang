package main

import (
	"github.com/josejuanmontiel/golang/testing/doer"
	"github.com/josejuanmontiel/golang/testing/user"
)

func main() {
	newDoerImpl := doer.NewDoerService()
	_ = newDoerImpl.DoSomething(1, "hola")

	newUser := user.NewUser(newDoerImpl)
	newUser.Use()
}
