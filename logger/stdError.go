package logger

import (
	"fmt"
)

type StdError struct {
	errorCount int
}

func (stdError *StdError) Log(location string, message string) {
	stdError.errorCount++
	fmt.Println(location + ": " + message)
}

func (stdError *StdError) ErrorsCount() int {
	return stdError.errorCount
}
