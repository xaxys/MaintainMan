package util

import (
	"bytes"
	"html/template"
	"math/rand"
	"regexp"
)

var (
	EmailRegex = regexp.MustCompile(`^[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+$`)
	PhoneRegex = regexp.MustCompile(`^[0-9]{3}[-\s\.]?[0-9]{4}([-\s\.]?[0-9]{4})?$`)
)

func RandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func ToUint[T ~int | ~int8 | ~int16 | ~int32 | ~int64](n T) uint {
	if n < 0 {
		return 0
	}
	return uint(n)

}

func NotEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func NotNil[T, U any](v *T, u *T) *T {
	if v == nil {
		return u
	}
	return v
}

func NilOrValue[T, U any](v *T, obj *U) *U {
	if v == nil {
		return nil
	}
	return obj

}

func NilOrLazyValue[T, U any](v *T, fn func(*T) *U) *U {
	if v == nil {
		return nil
	}
	return fn(v)

}

func NilOrBaseValue[T, U any](v *T, fn func(*T) U, def U) U {
	if v == nil {
		return def
	}
	return fn(v)

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
	}
	return b

}

func In[T comparable](key T, values ...T) bool {
	for _, v := range values {
		if key == v {
			return true
		}
	}
	return false
}

func Insert[T any](slice []T, elem T, index uint) []T {
	return append(slice[:index], append([]T{elem}, slice[index:]...)...)
}

func Move[T comparable](slice []T, elem T, to uint) []T {
	if int(to) >= len(slice) {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice[i], slice[to] = slice[to], slice[i]
			break
		}
	}
	return slice
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

func ProcessString(str string, vars interface{}) string {
	tmpl, err := template.New("").Parse(str)
	if err != nil {
		return str
	}
	buffer := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buffer, vars); err != nil {
		panic(err)
	}
	return buffer.String()
}
