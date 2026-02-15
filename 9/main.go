package main

import (
	"errors"
	"strings"
	"unicode"
)

func Unpack(s string) (string, error) {
	var builder strings.Builder
	var lastRune rune
	isEscaped := false

	for _, r := range s {
		if isEscaped {
			if lastRune != 0 {
				builder.WriteRune(lastRune)
			}
			lastRune = r
			isEscaped = false
			continue
		}

		if r == '\\' {
			isEscaped = true
			continue
		}

		if unicode.IsDigit(r) {
			if lastRune == 0 {
				return "", errors.New("invalid string: starts with digit or double digit")
			}
			repeatCount := int(r - '0')
			for i := 0; i < repeatCount; i++ {
				builder.WriteRune(lastRune)
			}
			lastRune = 0
			continue
		}

		// Если это обычный символ
		if lastRune != 0 {
			builder.WriteRune(lastRune)
		}
		lastRune = r
	}

	// Проверки после завершения цикла
	if isEscaped {
		return "", errors.New("invalid string: trailing backslash")
	}
	if lastRune != 0 {
		builder.WriteRune(lastRune)
	}

	return builder.String(), nil
}
