package user

import (
	"fmt"

	"github.com/josejuanmontiel/golang/testing/doer"
)

type User struct {
	Doer doer.Doer
}

func NewUser(d doer.Doer) *User {
	return &User{d}
}

func (u *User) Use() error {
	fmt.Println("Calling Use method")
	return u.Doer.DoSomething(123, "Hello GoMock")
}
