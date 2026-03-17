package main

import (
	"errors"
	"fmt"
	"unicode"
)

func CreateRUCFromDNI(dni string) (string, error) {
	if len(dni) != 8 {
		return "", errors.New("dni must have 8 digits")
	}
	for _, r := range dni {
		if !unicode.IsDigit(r) {
			return "", errors.New("dni must contain only digits")
		}
	}

	base := "10" + dni
	checkDigit := calculateRUCCheckDigit(base)

	return base + fmt.Sprintf("%d", checkDigit), nil
}

func calculateRUCCheckDigit(base string) int {
	weights := []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i, r := range base {
		digit := int(r - '0')
		sum += digit * weights[i]
	}

	remainder := sum % 11
	check := 11 - remainder

	switch check {
	case 10:
		return 0
	case 11:
		return 1
	default:
		return check
	}
}
func LeftPadZero(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return fmt.Sprintf("%0*s", length, s)
}
