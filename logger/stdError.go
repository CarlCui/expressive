package logger

import (
	"fmt"
)

type StdError int

func (stdError StdError) Log(location string, message string) {
	stdError++
	fmt.Println(location + ": " + message)
}
