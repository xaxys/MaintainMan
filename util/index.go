package util

import "strconv"

func ParseInt(b string, dft int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return dft
	} else {
		return id
	}
}
