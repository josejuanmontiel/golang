package doer

import (
)

type Doer interface {
    DoSomething(int, string) error
}
