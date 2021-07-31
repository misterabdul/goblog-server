package utils

import "fmt"

/**
 * Colored console ouput's text.
 * Reference from: https://golangbyexample.com/print-output-text-color-console/
 */

// Get console's color value from given name.
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

// Print given string to console with given color name.
//
// Available color name: white, green, yellow, red, blue, purple, and cyan.
func ConsolePrintColor(str string, color string) (int, error) {
	return fmt.Println(getColor(color), str)
}

// Print given string to console with white color.
func ConsolePrintWhite(str string) (int, error) {
	return ConsolePrintColor(str, "white")
}

// Print given string to console with yellow color.
func ConsolePrintYellow(str string) (int, error) {
	return ConsolePrintColor(str, "yellow")
}

// Print given string to console with green color.
func ConsolePrintGreen(str string) (int, error) {
	return ConsolePrintColor(str, "green")
}

// Print given string to console with red color.
func ConsolePrintRed(str string) (int, error) {
	return ConsolePrintColor(str, "red")
}
