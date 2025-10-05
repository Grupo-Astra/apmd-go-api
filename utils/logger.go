// Package utils fornece funções auxiliares e utilitários
// para toda a aplicação.
package utils

import (
	"fmt"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// LogInfo exibe uma mensagem informativa no console, com formatação ciano.
func LogInfo(msg string) {
	fmt.Println(colorCyan + "[INFO] " + colorReset + msg)
}

// LogSuccess exibe uma mensagem de sucesso no console, com formatação verde.
func LogSuccess(msg string) {
	fmt.Println(colorGreen + "[✔] " + colorReset + msg)
}

// LogWarn exibe uma mensagem de aviso no console, com formatação amarela.
func LogWarn(msg string) {
	fmt.Println(colorYellow + "[⚠] " + colorReset + msg)
}

// LogError exibe uma mensagem de erro no console, com formatação vermelha.
func LogError(msg string) {
	fmt.Println(colorRed + "[✖] " + colorReset + msg)
}

// LogSection exibe um título de seção no console,
// para organização visual dos logs.
func LogSection(title string) {
	fmt.Println(colorBlue + "=== " + title + " ===" + colorReset)
}

// LogFatal exibe uma mensagem de erro fatal e encerra a
// aplicação com status 1.
func LogFatal(msg string) {
	fmt.Println(colorRed + "[FATAL] " + colorReset + msg)
	os.Exit(1)
}
