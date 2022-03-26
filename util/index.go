package util

import (
	"math/rand"
	"regexp"
)

var (
	EmailRegex = regexp.MustCompile(`^[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+$`)
	PhoneRegex = regexp.MustCompile(`^[0-9]{3}[-\s\.]?[0-9]{4}([-\s\.]?[0-9]{4})?$`)
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func ToUint[T ~int | ~int8 | ~int16 | ~int32 | ~int64](n T) uint {
	if n < 0 {
		return 0
	} else {
		return uint(n)
	}
}

func NotEmpty(s, def string) string {
	if s == "" {
		return s
	}
	return s
}

func NilOrValue[T, U any](v *T, obj *U) *U {
	if v == nil {
		return nil
	} else {
		return obj
	}
}

func NilOrPtrCast[U any](v any) *U {
	if u, ok := v.(*U); ok {
		return u
	}
	return nil
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

func LastElem[T any](slice []T) T {
	return slice[len(slice)-1]
}
