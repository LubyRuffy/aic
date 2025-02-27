package color

import (
	"fmt"
	"github.com/fatih/color"
)

// Success prints success message in green
func Success(format string, a ...interface{}) {
	color.Green(format, a...)
}

// Info prints info message in blue
func Info(format string, a ...interface{}) {
	color.Blue(format, a...)
}

// Warning prints warning message in yellow
func Warning(format string, a ...interface{}) {
	color.Yellow(format, a...)
}

// Error prints error message in red
func Error(format string, a ...interface{}) {
	fmt.Print(color.RedString(format, a...))
}
