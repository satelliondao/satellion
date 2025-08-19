package term

import "fmt"

const clearLine = "\r\033[K"

// PrintInline clears the current line and writes the given message without a trailing newline.
func PrintInline(message string) {
	fmt.Print(clearLine)
	fmt.Print(message)
}

// PrintfInline clears the current line and writes the formatted message without a trailing newline.
func PrintfInline(format string, a ...any) {
	fmt.Print(clearLine)
	fmt.Printf(format, a...)
}

// Newline writes a newline, useful after inline updates before final messages.
func Newline() {
	fmt.Print("\n")
}


