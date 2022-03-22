package util

import (
	"math/rand"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func NotNil[T, U any](v *T, obj *U) *U {
	if v == nil {
		return nil
	} else {
		return obj
	}
}

func Tenary[T any](cond bool, a, b T) T {
	if cond {
		return a
	} else {
		return b
	}
}

func In[T comparable](key T, values ...T) bool {
	for _, v := range values {
		if key == v {
			return true
		}
	}
	return false
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
