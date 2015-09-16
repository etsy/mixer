package debug

import "fmt"

type Debug bool

func (d Debug) Printf(s string, a ...interface{}) {
	if d {
		fmt.Printf(s, a...)
	}
}
