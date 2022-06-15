package chainables_test

import (
	"errors"
	"fmt"

	"github.com/ojizero/hades/chainables"
)

func ExampleWith() {
	foo := func() (int, error) {
		return 0, errors.New("some error")
	}
	bar := func(i int) {
		fmt.Println("This won't be reached in case of error from foo since it is chained after it.")
		fmt.Printf("If it reached it'll receive the output of foo as an arg -> (%d) <-", i)
	}
	chainables.With(foo, bar)
}
