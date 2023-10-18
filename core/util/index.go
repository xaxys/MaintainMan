package util

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"reflect"
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

func CopySlice[T any](slice []T) []T {
	return append([]T(nil), slice...)
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
	return ProcessStringWithTPL(tmpl, str, vars)
}

func ProcessStringWithTPL(tpl *template.Template, str string, vars interface{}) string {
	tmpl, err := tpl.Parse(str)
	if err != nil {
		return str
	}
	buffer := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buffer, vars); err != nil {
		panic(err)
	}
	return buffer.String()
}

func FormatBytes(bytes uint64) string {
	if bytes>>40 > 100 {
		return fmt.Sprintf("%d TiB", bytes>>40)
	} else if bytes>>30 > 100 {
		return fmt.Sprintf("%d GiB", bytes>>30)
	} else if bytes>>20 > 100 {
		return fmt.Sprintf("%d MiB", bytes>>20)
	} else if bytes>>10 > 100 {
		return fmt.Sprintf("%d KiB", bytes>>10)
	} else {
		return fmt.Sprintf("%d B", bytes)
	}
}

// NotEmptyFieldName Get the json tag name of a struct field that is not empty
func NotEmptyFieldName(s any) (names []string) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ok := false
		switch field.Kind() {
		case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
			ok = field.Len() > 0
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ok = field.Int() != 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ok = field.Uint() != 0
		case reflect.Float32, reflect.Float64:
			ok = field.Float() != 0
		case reflect.Bool:
			ok = field.Bool()
		case reflect.Struct:
			ok = NotEmptyFieldName(field.Interface()) != nil
		case reflect.Ptr:
			ok = field.Elem().IsValid()
		}
		name := ""
		if !ok {
			continue
		}
		for _, tag := range []string{"json", "yaml", "toml", "xml", "bson", "url"} {
			name = t.Field(i).Tag.Get(tag)
			if name != "" {
				break
			}
		}
		if name == "" {
			name = t.Field(i).Name
		}
		names = append(names, name)
	}
	return names
}
