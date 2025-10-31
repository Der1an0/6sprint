package service

import (
	"fmt"
	"strings"

	"github.com/Der1an0/6sprint/pkg/morse" // импорт пакета morse
)

// Определяем тип текста и конвертируем его
func Convert(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("empty line")
	}

	if isMorseCode(input) {
		// Конвертируем Морзе в текст
		return morse.ToText(input), nil
	} else {
		// Конвертируем текст в Морзе
		return morse.ToMorse(input), nil
	}
}

// Проверяем, является ли строка кодом Морзе
func isMorseCode(input string) bool {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return false
	}
	// Код Морзе состоит только из точек, тире, пробелов и слэшей
	for _, char := range trimmed {
		if char != '.' && char != '-' && char != ' ' && char != '/' && char != '\n' && char != '\r' && char != '\t' {
			return false
		}
	}
	// Дополнительная проверка: если есть хотя бы одна точка или тире - вероятно это Морзе
	hasMorseChars := strings.ContainsAny(trimmed, ".-")

	return hasMorseChars
}
