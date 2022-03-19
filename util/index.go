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

func Remove[T comparable](slice []T, elems ...T) []T {
	for _, e := range elems {
		for i, v := range slice {
			if v == e {
				return append(slice[:i], slice[i+1:]...)
			}
		}
	}
	return slice
}

func RemoveByRef[T any](slice []T, elems ...*T) []T {
	for _, e := range elems {
		for i := range slice {
			if &slice[i] == e {
				return append(slice[:i], slice[i+1:]...)
			}
		}
	}
	return slice
}

func TransSlice[T, U any](s []T, trans func(T) U) (us []U) {
	for _, t := range s {
		us = append(us, trans(t))
	}
	return
}
