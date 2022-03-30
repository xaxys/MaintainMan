package util

import (
	"sync"
	"sync/atomic"
)

type CoPtrMap[K, V any] struct {
	data sync.Map
}

func NewCoMap[K, V any]() *CoPtrMap[K, *V] {
	return &CoPtrMap[K, *V]{}
}

func (c *CoPtrMap[K, V]) Get(k K) *V {
	if v, ok := c.data.Load(k); ok {
		return v.(*V)
	}
	return nil
}

func (c *CoPtrMap[K, V]) Has(k K) bool {
	_, ok := c.data.Load(k)
	return ok
}

func (c *CoPtrMap[K, V]) Set(k K, v *V) {
	c.data.Store(k, v)
}

func (c *CoPtrMap[K, V]) Delete(k K) {
	c.data.Delete(k)
}

func (c *CoPtrMap[K, V]) LoadAndDelete(k K) *V {
	if value, ok := c.data.LoadAndDelete(k); ok {
		return value.(*V)
	}
	return nil
}

func (c *CoPtrMap[K, V]) Range(f func(k K, v *V) error) (err error) {
	c.data.Range(func(k, v any) bool {
		if err = f(k.(K), v.(*V)); err != nil {
			return false
		}
		return true
	})
	return
}

type AtomValue[T any] struct {
	value atomic.Value
}

func NewAtomValue[T any](v T) *AtomValue[T] {
	a := &AtomValue[T]{}
	a.value.Store(v)
	return a
}

func (a *AtomValue[T]) Get() T {
	return a.value.Load().(T)
}

func (a *AtomValue[T]) Set(v T) {
	a.value.Store(v)
}
