package log

import (
	"fmt"
	"os"
)

func Print(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Log(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Error(format string, a ...any) {
	fmt.Printf("error: "+format+"\n", a...)
}

func Fatal(format string, a ...any) {
	Error(format, a...)
	os.Exit(1)
}
