package util

import (
	"regexp"
	"strconv"
)

var (
	EmailRegex = regexp.MustCompile(`^[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+$`)
	PhoneRegex = regexp.MustCompile(`^[0-9]{3}[-\s\.]?[0-9]{4}([-\s\.]?[0-9]{4})?$`)
)

func ParseInt(b string, dft int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return dft
	} else {
		return id
	}
}
