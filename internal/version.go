package internal

import (
	"fmt"
)

var VERSION string

func PrintVersion() {
	fmt.Printf("%v %v\n", NAME, VERSION)
}
