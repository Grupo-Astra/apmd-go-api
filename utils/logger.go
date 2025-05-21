package utils

import "fmt"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

func LogInfo(msg string) {
	fmt.Println(colorCyan + "[INFO] " + colorReset + msg)
}

func LogSuccess(msg string) {
	fmt.Println(colorGreen + "[✔] " + colorReset + msg)
}

func LogWarn(msg string) {
	fmt.Println(colorYellow + "[⚠] " + colorReset + msg)
}

func LogError(msg string) {
	fmt.Println(colorRed + "[✖] " + colorReset + msg)
}

func LogSection(title string) {
	fmt.Println(colorBlue + "=== " + title + " ===" + colorReset)
}
