package main

import (
	"errors"
	"regexp"
)

func IsValidRuc(ruc string) error {
	if match, err := regexp.MatchString("^[0-9]{11}$", ruc); err != nil || !match {
		return errors.New("RUC invalido")
	}

	sa := []rune(ruc)
	sum := 0
	x := 6
	d := 0
	for i := 0; i < 10; i++ {
		if i == 4 {
			x = 8
		}
		d = int(sa[i] - '0')
		x = x - 1
		sum = sum + d*x
	}
	r := sum % 11
	r = 11 - r
	if r >= 10 {
		r = r - 10
	}

	if r == int(sa[10]-'0') {
		return nil
	}
	return errors.New("RUC invalido")
}
