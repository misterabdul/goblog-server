package utils

import "fmt"

/**
 * Colored console ouput's text.
 * Reference from: https://golangbyexample.com/print-output-text-color-console/
 */

func getColor(name string) string {
	switch {
	case name == "reset":
		return "\033[0m"
	case name == "red":
		return "\033[31m"
	case name == "green":
		return "\033[32m"
	case name == "yellow":
		return "\033[33m"
	case name == "blue":
		return "\033[34m"
	case name == "purple":
		return "\033[35m"
	case name == "cyan":
		return "\033[36m"
	case name == "white":
		return "\033[37m"
	}
	return "\033[37m"
}

func ConsolePrintColor(str string, color string) (int, error) {
	return fmt.Println(getColor(color), str)
}

func ConsolePrintWhite(str string) (int, error) {
	return ConsolePrintColor(str, "white")
}

func ConsolePrintYellow(str string) (int, error) {
	return ConsolePrintColor(str, "yellow")
}

func ConsolePrintGreen(str string) (int, error) {
	return ConsolePrintColor(str, "green")
}

func ConsolePrintRed(str string) (int, error) {
	return ConsolePrintColor(str, "red")
}
