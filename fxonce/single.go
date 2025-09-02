package fxonce

import "sync"

// Do 单列模式
func Do[T any](fn func() T) func() T {
	var (
		once sync.Once
		s    T
	)
	return func() T {
		once.Do(func() {
			s = fn()
		})
		return s
	}
}

// DoWithParam 单列带参数模式 泛型函数
func DoWithParam[T any, P any](fn func(P) T) func(P) T {
	var (
		once sync.Once
		s    T
	)
	return func(param P) T {
		once.Do(func() {
			s = fn(param)
		})
		return s
	}
}
